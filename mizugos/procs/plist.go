package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
)

// NewPList 建立plist處理器
func NewPList() *PList {
	return &PList{
		Procmgr: NewProcmgr(),
	}
}

// PList plist處理器, 封包結構使用PListMsg, 由於使用到des-cbc加密, 安全性較高, 適合用來傳送一般封包, 使用時需要設定傳送函式, 密鑰以及初始向量;
// 由於採用複數訊息設計, 因此使用plist上下文與使用者溝通(json/proto處理器則使用訊息結構與使用者溝通);
// 當使用者需要回應訊息給客戶端時, 把回應訊息新增到plist上下文中, 在處理的最後會一次性傳送所有回應訊息給客戶端
//   - 訊息內容: support/proto/mizugo/plistmsg.proto
//   - 封包編碼: protobuf編碼成位元陣列, 再通過des加密
//   - 封包解碼: des解密, 再通過protobuf解碼成訊息結構
type PList struct {
	*Procmgr                        // 管理器
	send     utils.SyncAttr[Send]   // 傳送函式
	key      utils.SyncAttr[[]byte] // 密鑰
	iv       utils.SyncAttr[[]byte] // 初始向量
}

// Send 傳送函式類型
type Send func(message any)

// Send 設定傳送函式
func (this *PList) Send(send Send) *PList {
	this.send.Set(send)
	return this
}

// Key 設定密鑰
func (this *PList) Key(key []byte) *PList {
	this.key.Set(key)
	return this
}

// KeyStr 設定密鑰
func (this *PList) KeyStr(key string) *PList {
	return this.Key([]byte(key))
}

// IV 設定初始向量
func (this *PList) IV(iv []byte) *PList {
	this.iv.Set(iv)
	return this
}

// IVStr 設定初始向量
func (this *PList) IVStr(iv string) *PList {
	return this.IV([]byte(iv))
}

// Encode 封包編碼
func (this *PList) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[msgs.PListMsg](input)

	if err != nil {
		return nil, fmt.Errorf("plist encode: %w", err)
	} // if

	bytes, err := proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("plist encode: %w", err)
	} // if

	output, err = cryptos.DesCBCEncrypt(cryptos.PaddingPKCS7, this.key.Get(), this.iv.Get(), bytes)

	if err != nil {
		return nil, fmt.Errorf("plist encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *PList) Decode(input []byte) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("plist decode: input nil")
	} // if

	bytes, err := cryptos.DesCBCDecrypt(cryptos.PaddingPKCS7, this.key.Get(), this.iv.Get(), input)

	if err != nil {
		return nil, fmt.Errorf("plist decode: %w", err)
	} // if

	message := &msgs.PListMsg{}

	if err = proto.Unmarshal(bytes, message); err != nil {
		return nil, fmt.Errorf("plist decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *PList) Process(input any) error {
	message, err := utils.CastPointer[msgs.PListMsg](input)

	if err != nil {
		return fmt.Errorf("plist process: %w", err)
	} // if

	context := &PListContext{
		request: message.Messages,
	}

	for context.next() {
		messageID := context.messageID()
		process := this.Get(messageID)

		if process == nil {
			return fmt.Errorf("plist process: not found: %v", messageID)
		} // if

		process(context)
	} // if

	send := this.send.Get()

	if send == nil {
		return fmt.Errorf("plist process: send nil")
	} // if

	result, err := PListMarshal(context)

	if err != nil {
		return fmt.Errorf("plist process: %w", err)
	} // if

	send(result)
	return nil
}

// PListContext plist上下文資料
type PListContext struct {
	request []*msgs.PListUnit // 要求訊息列表
	respond []*msgs.PListUnit // 回應訊息列表
	currmsg *msgs.PListUnit   // 當前訊息
	currpos int               // 當前位置
}

// next 移動到下個訊息
func (this *PListContext) next() bool {
	if this.currpos < len(this.request) {
		this.currmsg = this.request[this.currpos]
		this.currpos++
		return true
	} // if

	return false
}

// messageID 取得當前訊息編號
func (this *PListContext) messageID() MessageID {
	if this.currmsg != nil {
		return this.currmsg.MessageID
	} // if

	return 0
}

// message 取得當前訊息
func (this *PListContext) message() *anypb.Any {
	if this.currmsg != nil {
		return this.currmsg.Message
	} // if

	return nil
}

// AddRespond 新增回應訊息
func (this *PListContext) AddRespond(messageID MessageID, input proto.Message) error {
	message, err := anypb.New(input)

	if err != nil {
		return fmt.Errorf("plist context addrespond: %w", err)
	} // if

	this.respond = append(this.respond, &msgs.PListUnit{
		MessageID: messageID,
		Message:   message,
	})
	return nil
}

// PListMarshal 序列化plist訊息
func PListMarshal(input *PListContext) (output *msgs.PListMsg, err error) {
	if input == nil {
		return nil, fmt.Errorf("plist marshal: input nil")
	} // if

	return &msgs.PListMsg{
		Messages: input.respond,
	}, nil
}

// PListUnmarshal 反序列化plist訊息
func PListUnmarshal[T any](input *PListContext) (messageID MessageID, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("plist unmarshal: input nil")
	} // if

	message := input.message()

	if message == nil {
		return 0, nil, fmt.Errorf("plist unmarshal: message nil")
	} // if

	temp, err := message.UnmarshalNew()

	if err != nil {
		return 0, nil, fmt.Errorf("plist unmarshal: %w", err)
	} // if

	output, ok := temp.(any).(*T)

	if ok == false {
		return 0, nil, fmt.Errorf("plist unmarshal: cast failed")
	} // if

	return input.messageID(), output, nil
}
