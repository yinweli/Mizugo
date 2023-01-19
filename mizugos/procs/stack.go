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

	context := &StackContext{
		request: message.Messages,
	}

	for context.next() {
		messageID := context.messageID()
		process := this.Get(messageID)

		if process == nil {
			return fmt.Errorf("stack process: not found: %v", messageID)
		} // if

		process(context)
	} // if

	send := this.send.Get()

	if send == nil {
		return fmt.Errorf("stack process: send nil")
	} // if

	result, err := StackMarshal(context)

	if err != nil {
		return fmt.Errorf("stack process: %w", err)
	} // if

	send(result)
	return nil
}

// StackContext 堆棧上下文資料
type StackContext struct {
	request []*msgs.StackUnit // 要求訊息列表
	respond []*msgs.StackUnit // 回應訊息列表
	currmsg *msgs.StackUnit   // 當前訊息
	currpos int               // 當前位置
}

// next 移動到下個訊息
func (this *StackContext) next() bool {
	if this.currpos < len(this.request) {
		this.currmsg = this.request[this.currpos]
		this.currpos++
		return true
	} // if

	return false
}

// messageID 取得當前訊息編號
func (this *StackContext) messageID() MessageID {
	if this.currmsg != nil {
		return this.currmsg.MessageID
	} // if

	return 0
}

// message 取得當前訊息
func (this *StackContext) message() *anypb.Any {
	if this.currmsg != nil {
		return this.currmsg.Message
	} // if

	return nil
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

// StackMarshal 序列化堆棧訊息
func StackMarshal(input *StackContext) (output *msgs.StackMsg, err error) {
	if input == nil {
		return nil, fmt.Errorf("stack marshal: input nil")
	} // if

	return &msgs.StackMsg{
		Messages: input.respond,
	}, nil
}

// StackUnmarshal 反序列化堆棧訊息
func StackUnmarshal[T any](input *StackContext) (messageID MessageID, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("stack unmarshal: input nil")
	} // if

	message := input.message()

	if message == nil {
		return 0, nil, fmt.Errorf("stack unmarshal: message nil")
	} // if

	temp, err := message.UnmarshalNew()

	if err != nil {
		return 0, nil, fmt.Errorf("stack unmarshal: %w", err)
	} // if

	output, ok := temp.(any).(*T)

	if ok == false {
		return 0, nil, fmt.Errorf("stack unmarshal: cast failed")
	} // if

	return input.messageID(), output, nil
}
