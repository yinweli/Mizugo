package nets

import (
	"net"
)

const ( // 網路事件名稱
	EventRecv = "recv" // 接收訊息事件, 當接收訊息後觸發, 參數是訊息物件
	EventSend = "send" // 傳送訊息事件, 當傳送訊息後觸發, 參數是訊息物件
)

// Connecter 連接介面
type Connecter interface {
	// Connect 啟動連接
	Connect(bind Bind, unbind Unbind, wrong Wrong)

	// Address 取得位址
	Address() string
}

// Listener 接聽介面
type Listener interface {
	// Listen 啟動接聽
	Listen(bind Bind, unbind Unbind, wrong Wrong)

	// Stop 停止接聽
	Stop() error

	// Address 取得位址
	Address() string
}

// Sessioner 會話介面
type Sessioner interface {
	// Start 啟動會話, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止會話; 當由連接器/接聽器獲得會話器之後, 需要啟動會話才可以傳送或接收封包
	Start(bind Bind, unbind Unbind, wrong Wrong)

	// Stop 停止會話, 不會等待會話內部循環結束
	Stop()

	// StopWait 停止會話, 會等待會話內部循環結束
	StopWait()

	// Send 傳送封包
	Send(message any)

	// RemoteAddr 取得遠端位址
	RemoteAddr() net.Addr

	// LocalAddr 取得本地位址
	LocalAddr() net.Addr

	// SetOwner 設定擁有者
	SetOwner(owner any)

	// GetOwner 取得擁有者
	GetOwner() any
}

// Bind 綁定處理函式類型
type Bind func(session Sessioner) *Bundle

// Do 執行處理
func (this Bind) Do(session Sessioner) *Bundle {
	if this != nil {
		return this(session)
	} // if

	return nil
}

// Unbind 解綁處理函式類型
type Unbind func(session Sessioner)

// Do 執行處理
func (this Unbind) Do(session Sessioner) {
	if this != nil {
		this(session)
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

// Bundle 綁定資料
type Bundle struct {
	Encode
	Decode
	Publish
}

// Encode 封包編碼處理函式類型, 用在傳送封包時
type Encode func(message any) (packet []byte, err error)

// Decode 封包解碼處理函式類型, 用在接收封包時
type Decode func(packet []byte) (message any, err error)

// Publish 發布事件處理函式類型
type Publish func(name string, param any)

// Do 執行處理
func (this Publish) Do(name string, param any) {
	if this != nil {
		this(name, param)
	} // if
}
