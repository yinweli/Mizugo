package nets

import (
	"net"
)

// nets 提供抽象化的網路介面: Connecter / Listener / Sessioner / Codec 與一組回呼函式型別 Bind / Unbind / Publish / Wrong
// 讓呼叫端能以一致的方式管理連線, 接聽與會話生命週期
//
// # 整體使用流程建議
//
//   1) 由 Listener(伺服器) 或 Connecter(客戶端) 取得 Sessioner
//   2) 在啟動會話前, 先以 Sessioner 設定:
//      - SetCodec: 編解碼鏈(Encode: 前→後; Decode: 後→前)
//      - SetPublish: 事件通知 EventStart / EventStop / EventRecv / EventSend
//      - SetWrong: 錯誤通知, 集中處理發生在底層收發/解碼/編碼等錯誤
//      - SetPacketSize: 若需覆蓋預設封包長度上限(單位: 位元組)
//      - SetOwner: 關聯上層物件(例如: 對應的實體, 玩家或連線上下文)
//   3) 呼叫 Sessioner.Start(bind, unbind) 啟動收發循環
//   4) 收到事件與錯誤後, 依需求處理
//   5) 停止會話可呼叫 Sessioner.Stop / Sessioner.StopWait
//   6) 停止接聽/連線由各自介面提供
//
// # 事件與回呼的語義(慣例, 視實作而定)
//
// - Bind(session) bool:
//   在會話啟動「前」被呼叫, 用於初始化(設 Owner, Codec, 事件/錯誤處理等)
//   回傳 true 表示允許啟動; false 表示拒絕(通常實作端會立即關閉連線)
// - Unbind(session):
//   在會話停止「後」被呼叫, 用於釋放資源
// - Publish(name, param):
//   用於對外發佈會話事件; name 通常為 EventStart / EventStop / EventRecv / EventSend
//   param 會依事件不同而異(例如: EventStart / EventStop 常傳 session, EventRecv / EventSend 常傳訊息)
// - Wrong(err):
//   收發循環或編解碼出錯時呼叫; 建議在此集中記錄/計數/告警與回收資源
//
// # 關於編解碼鏈(Codec)
//
// - SetCodec(c1, c2, c3):
//   Encode(傳送)時的順序: c1 → c2 → c3
//   Decode(接收)時的順序: c3 → c2 → c1 (反序)
// - Codec 輸入/輸出型別為 any, 因此鏈中的每個編解碼器都需清楚定義其期望與產出型別
//   Encode 產出為 []byte
//   Decode 輸入為 []byte
//
// # 封包長度與標頭
//
// - HeaderSize 為 4, 表示實作上將以 4-byte 無號整數記錄封包長度
// - PacketSize 為單一封包內容(payload)的預設長度上限(不含標頭)
// - 具體大小端(endianness)依會話實作而定
// - 若搭配本套件 TCP 實作, 預設採 LittleEndian, 若跨系統/語言溝通, 請務必確保雙方對標頭大小端一致

const ( // 網路定義
	// HeaderSize 為封包標頭長度(位元組), 用於記錄「payload 長度」
	// 一般對應 4-byte 無號整數; 端序請與實際 Sessioner 實作一致(搭配 TCP 預設 LittleEndian)
	HeaderSize = 4

	// PacketSize 為「單一封包內容」的預設長度上限(位元組), 不含標頭
	// 若需要更大封包, 請在會話層以 Sessioner.SetPacketSize 覆蓋
	PacketSize = int(^uint16(0))

	// ChannelSize 為內部訊息通道的緩衝大小若應用在高吞吐環境, 通道爆滿時, 傳送端可能被阻塞
	ChannelSize = 1000
)

const ( // 網路事件名稱
	// EventStart 會話已成功啟動並進入收發循環
	// 參數: Sessioner
	EventStart = "start"

	// EventStop 會話已停止並完成收尾
	// 參數: Sessioner
	EventStop = "stop"

	// EventRecv 會話層完成解碼並收到一則訊息
	// 參數: 收到的訊息(型別由上層協議定義)
	EventRecv = "recv"

	// EventSend 會話層完成編碼並成功送出一則訊息
	// 參數: 送出的訊息(通常為編碼前的上層物件)
	EventSend = "send"
)

// Connecter 連接介面
type Connecter interface {
	// Connect 啟動連接
	//
	// 流程:
	//   - 以實作(例如 TCP)撥出連線; 失敗時呼叫 wrong(err)
	//   - 連線成功後建立會話(Sessioner)並呼叫 Sessioner.Start(bind, unbind)
	//
	// Connect 的錯誤回報慣例透過 wrong callback 傳遞
	Connect(bind Bind, unbind Unbind, wrong Wrong)

	// Address 取得連接位址字串
	Address() string
}

// Listener 接聽介面
type Listener interface {
	// Listen 啟動接聽
	//
	// 流程:
	//   - 以實作(例如 TCP)綁定位址並開始 Accept; 失敗時呼叫 wrong(err)
	//   - 每次 Accept 連線後建立會話(Sessioner)並呼叫 Sessioner.Start(bind, unbind)
	//
	// Listen 一般為非阻塞
	Listen(bind Bind, unbind Unbind, wrong Wrong)

	// Stop 停止接聽
	Stop() error

	// Address 取得接聽位址字串
	Address() string
}

// Sessioner 會話介面
type Sessioner interface {
	// Start 啟動會話
	//
	// 注意:
	//   - 若非使用多執行緒, 呼叫端可能被阻塞直到會話停止
	//   - 需於 Start 前完成必要設定(如 SetCodec / SetPublish / SetWrong / SetPacketSize / SetOwner)
	//   - 啟動成功後, 會進入內部收發循環, 並在適當時機觸發以下事件:
	//     Bind(session) -> EventStart -> (收發循環) -> EventStop -> Unbind(session)
	//
	// Bind / Unbind:
	//   - Bind 回傳 false 表示拒絕啟動(通常直接關閉連線)
	//   - Unbind 在會話結束後被呼叫, 用於釋放資源
	Start(bind Bind, unbind Unbind)

	// Stop 停止會話, 不等待收發循環自然結束
	// 若需要等待完整收尾, 請改用 StopWait
	Stop()

	// StopWait 停止會話並等待內部收發循環結束(阻塞直到關閉完成)
	StopWait()

	// SetPublish 設定事件發布處理器
	SetPublish(publish ...Publish)

	// SetWrong 設定錯誤處理器
	SetWrong(wrong ...Wrong)

	// SetCodec 設定編碼/解碼鏈
	//
	// 執行順序:
	//   - 編碼(Encode, 傳送): 按傳入順序由前到後
	//   - 解碼(Decode, 接收): 按傳入順序由後到前(反序)
	//
	// 介面型別為 any → any, 請確保串接的編解碼器能夠正確承接型別
	//
	// 最終 Encode 產物為 []byte; Decode 的初始輸入為 []byte
	SetCodec(codec ...Codec)

	// SetOwner 設定會話的「擁有者(實體)」關聯
	SetOwner(owner any)

	// GetOwner 取得 SetOwner 設定的擁有者
	GetOwner() any

	// SetPacketSize 設定單一封包內容最大長度上限(位元組, 不含標頭)
	SetPacketSize(size int)

	// Send 傳送訊息(型別由上層協議決定, 將經由編碼鏈轉為位元組流)
	Send(message any)

	// RemoteAddr 取得遠端位址
	RemoteAddr() net.Addr

	// LocalAddr 取得本地位址
	LocalAddr() net.Addr
}

// Codec 編碼/解碼介面
type Codec interface {
	// Encode 編碼處理:
	//   - 將上層訊息逐步轉換為 []byte
	//   - 若串接多個 Codec，會按 Sessioner.SetCodec 的傳入順序由前到後執行
	Encode(input any) (output any, err error)

	// Decode 解碼處理:
	//   - 將 []byte 逐步還原為上層訊息
	//   - 若串接多個 Codec，會按 Sessioner.SetCodec 的傳入順序由後到前執行
	Decode(input any) (output any, err error)
}

// Bind 綁定處理函式類型
type Bind func(session Sessioner) bool

// Do 執行處理
func (this Bind) Do(session Sessioner) bool {
	if this != nil {
		return this(session)
	} // if

	return true
}

// Unbind 解綁處理函式類型
type Unbind func(session Sessioner)

// Do 執行處理
func (this Unbind) Do(session Sessioner) {
	if this != nil {
		this(session)
	} // if
}

// Publish 發布事件處理函式類型
type Publish func(name string, param any)

// Do 執行處理
func (this Publish) Do(name string, param any) {
	if this != nil {
		this(name, param)
	} // if
}

// Wrong 錯誤處理函式類型
type Wrong func(err error)

// Do 執行處理
func (this Wrong) Do(err error) {
	if this != nil {
		this(err)
	} // if
}
