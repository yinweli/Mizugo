package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
)

// 堆棧處理器, 封包結構使用StackMsg
// 由於使用到des加密, 安全性較高, 適合用來傳送一般封包, 使用時需要設定傳送函式以及密鑰
// 由於採用複數訊息設計, 因此使用堆棧上下文與使用者溝通(json/proto處理器則使用訊息結構與使用者溝通)
// 當使用者需要回應訊息給客戶端時, 把回應訊息新增到堆棧上下文中, 在處理的最後會一次性傳送所有回應訊息給客戶端
// 訊息內容: support/proto/mizugo/stackmsg.proto
// 封包編碼: protobuf編碼成位元陣列, 再通過des加密
// 封包解碼: des解密, 再通過protobuf解碼成訊息結構

// NewStack 建立堆棧處理器
func NewStack() *Stack {
	return &Stack{
		procmgr: newProcmgr(),
	}
}

// Stack 堆棧處理器
type Stack struct {
	*procmgr                        // 管理器
	send     utils.SyncAttr[Send]   // 傳送函式
	key      utils.SyncAttr[[]byte] // 密鑰
}

// Send 傳送函式類型
type Send func(message any)

// Send 設定傳送函式
func (this *Stack) Send(send Send) *Stack {
	this.send.Set(send)
	return this
}

// Key 設定密鑰
func (this *Stack) Key(key []byte) *Stack {
	this.key.Set(key)
	return this
}

// Encode 封包編碼
func (this *Stack) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[msgs.StackMsg](input)

	if err != nil {
		return nil, fmt.Errorf("stack encode: %w", err)
	} // if

	bytes, err := proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("stack encode: %w", err)
	} // if

	output, err = utils.DesEncrypt(this.key.Get(), bytes)

	if err != nil {
		return nil, fmt.Errorf("stack encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Stack) Decode(input []byte) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("stack decode: input nil")
	} // if

	bytes, err := utils.DesDecrypt(this.key.Get(), input)

	if err != nil {
		return nil, fmt.Errorf("stack decode: %w", err)
	} // if

	message := &msgs.StackMsg{}

	if err = proto.Unmarshal(bytes, message); err != nil {
		return nil, fmt.Errorf("stack decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Stack) Process(input any) error {
	message, err := utils.CastPointer[msgs.StackMsg](input)

	if err != nil {
		return fmt.Errorf("stack process: %w", err)
	} // if

	stackContext := &StackContext{}
	stackContext.initialize(message)

	for stackContext.next() {
		messageID := stackContext.messageID()
		process := this.Get(messageID)

		if process == nil {
			return fmt.Errorf("stack process: not found: %v", messageID)
		} // if

		process(stackContext)
	} // if

	if send := this.send.Get(); send != nil {
		send(stackContext.result())
	} // if

	return nil
}

// StackContext 堆棧上下文資料
type StackContext struct {
	request []*msgs.StackUnit // 要求訊息列表
	respond []*msgs.StackUnit // 回應訊息列表
	current int               // 當前位置
	message *msgs.StackUnit   // 當前訊息
}

// initialize 初始化處理
func (this *StackContext) initialize(message *msgs.StackMsg) {
	this.request = message.Messages
}

// next 移動到下個訊息
func (this *StackContext) next() bool {
	if this.current < len(this.request) {
		this.message = this.request[this.current]
		this.current++
		return true
	} // if

	return false
}

// messageID 取得當前訊息編號
func (this *StackContext) messageID() MessageID {
	if this.message != nil {
		return this.message.MessageID
	} // if

	return 0
}

// Unmarshal 反序列化當前訊息
func (this *StackContext) Unmarshal() (messageID MessageID, output proto.Message, err error) {
	if this.message == nil {
		return 0, nil, fmt.Errorf("stackcontext unmarshal: message nil")
	} // if

	if output, err = this.message.Message.UnmarshalNew(); err != nil {
		return 0, nil, fmt.Errorf("stackcontext unmarshal: %w", err)
	} // if

	return this.message.MessageID, output, nil
}

// bundle 取得結果訊息
func (this *StackContext) result() *msgs.StackMsg {
	return &msgs.StackMsg{
		Messages: this.respond,
	}
}

// AddRespond 新增回應訊息
func (this *StackContext) AddRespond(messageID MessageID, input proto.Message) error {
	message, err := anypb.New(input)

	if err != nil {
		return fmt.Errorf("stackcontext addrespond: %w", err)
	} // if

	this.respond = append(this.respond, &msgs.StackUnit{
		MessageID: messageID,
		Message:   message,
	})
	return nil
}
