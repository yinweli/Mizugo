package procs

import (
	"fmt"
	"strings"
	"unsafe"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
)

// NewRaven 建立 Raven 處理器
func NewRaven() *Raven {
	return &Raven{
		Procmgr: NewProcmgr(),
	}
}

// Raven 伺服器處理器
//
// 訊息結構基於 msgs.RavenS / msgs.RavenC
//   - messageID: 訊息編號
//   - header: 標頭訊息, 內容為任意 Proto 訊息
//   - request: 請求訊息, 內容為任意 Proto 訊息
//   - respond: 回應列表, 每筆皆為任意 Proto 訊息
//
// 對應的訊息定義檔:
//   - Go: support/proto-mizugo/msg-go/msgs/raven.pb.go
//   - C#: support/proto-mizugo/msg-cs/msgs/Raven.cs
//   - Proto 定義: support/proto-mizugo/raven.proto
//
// 職責:
//   - Encode: 將 *msgs.RavenC 物件轉為 []byte 以利傳輸
//   - Decode: 將 []byte 解碼回 *msgs.RavenS
//   - Process: 根據 messageID 呼叫對應的處理函式
//
// 另外提供以下工具
//   - RavenSBuilder / RavenSParser / RavenSData: 建立與解析伺服器訊息
//   - RavenIsMessageID / RavenIsErrID / RavenHeader / RavenRequest / RavenRespondAt / RavenRespondFind: 常用查詢工具
type Raven struct {
	*Procmgr // 管理器
}

// Encode 訊息編碼
//
// 輸入必須是 *msgs.RavenC, 輸出為 []byte
func (this *Raven) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("raven encode: input nil")
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return nil, fmt.Errorf("raven encode: input not *msgs.RavenC")
	} // if

	output, err = proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("raven encode: %w", err)
	} // if

	return output, nil
}

// Decode 訊息解碼
//
// 輸入必須是 []byte, 輸出為 *msgs.RavenS
func (this *Raven) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("raven decode: input nil")
	} // if

	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("raven decode: input not []byte")
	} // if

	message := &msgs.RavenS{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("raven decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
//
// 輸入必須是 *msgs.RavenS, 會依據 messageID 找出並執行已註冊的處理函式, 若找不到對應的處理函式則回傳錯誤
func (this *Raven) Process(input any) error {
	if input == nil {
		return fmt.Errorf("raven process: input nil")
	} // if

	message, ok := input.(*msgs.RavenS)

	if ok == false {
		return fmt.Errorf("raven process: input not *msgs.RavenS")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("raven process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// NewRavenClient 建立 Raven 客戶端處理器
func NewRavenClient() *RavenClient {
	return &RavenClient{
		Procmgr: NewProcmgr(),
	}
}

// RavenClient Raven 客戶端處理器
//
// 訊息結構基於 msgs.RavenS / msgs.RavenC
//   - messageID: 訊息編號
//   - header: 標頭訊息, 內容為任意 Proto 訊息
//   - request: 請求訊息, 內容為任意 Proto 訊息
//   - respond: 回應列表, 每筆皆為任意 Proto 訊息
//
// 對應的訊息定義檔:
//   - Go: support/proto-mizugo/msg-go/msgs/raven.pb.go
//   - C#: support/proto-mizugo/msg-cs/msgs/Raven.cs
//   - Proto 定義: support/proto-mizugo/raven.proto
//
// 職責:
//   - Encode: 將 *msgs.RavenS 物件轉為 []byte 以利傳輸
//   - Decode: 將 []byte 解碼回 *msgs.RavenC
//   - Process: 根據 messageID 呼叫對應的處理函式
//
// 另外提供以下工具
//   - RavenCBuilder / RavenCParser / RavenCData: 建立與解析客戶端訊息
//   - RavenIsMessageID / RavenIsErrID / RavenHeader / RavenRequest / RavenRespondAt / RavenRespondFind: 常用查詢工具
type RavenClient struct {
	*Procmgr // 管理器
}

// Encode 訊息編碼
//
// 輸入必須是 *msgs.RavenS, 輸出為 []byte
func (this *RavenClient) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("raven client encode: input nil")
	} // if

	message, ok := input.(*msgs.RavenS)

	if ok == false {
		return nil, fmt.Errorf("raven client encode: input not *msgs.RavenS")
	} // if

	output, err = proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("raven client encode: %w", err)
	} // if

	return output, nil
}

// Decode 訊息解碼
//
// 輸入必須是 []byte, 輸出為 *msgs.RavenC
func (this *RavenClient) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("raven client decode: input nil")
	} // if

	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("raven client decode: input not []byte")
	} // if

	message := &msgs.RavenC{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("raven client decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
//
// 輸入必須是 *msgs.RavenC, 會依據 messageID 找出並執行已註冊的處理函式, 若找不到對應的處理函式則回傳錯誤
func (this *RavenClient) Process(input any) error {
	if input == nil {
		return fmt.Errorf("raven client process: input nil")
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return fmt.Errorf("raven client process: input not *msgs.RavenC")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("raven client process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// RavenSBuilder 建立 *msgs.RavenS
func RavenSBuilder(messageID int32, header, request proto.Message) (output any, err error) {
	if header == nil {
		return nil, fmt.Errorf("ravenSBuilder: header nil")
	} // if

	if request == nil {
		return nil, fmt.Errorf("ravenSBuilder: request nil")
	} // if

	message := &msgs.RavenS{
		MessageID: messageID,
	}

	if message.Header, err = anypb.New(header); err != nil {
		return nil, fmt.Errorf("ravenSBuilder: header: %w", err)
	} // if

	if message.Request, err = anypb.New(request); err != nil {
		return nil, fmt.Errorf("ravenSBuilder: request: %w", err)
	} // if

	return message, nil
}

// RavenSParser 解析 *msgs.RavenS → RavenSData
func RavenSParser[H, Q proto.Message](input any) (output *RavenSData[H, Q], err error) {
	if input == nil {
		return nil, fmt.Errorf("ravenSParser: input nil")
	} // if

	message, ok := input.(*msgs.RavenS)

	if ok == false {
		return nil, fmt.Errorf("ravenSParser: input not *msgs.RavenS")
	} // if

	output = &RavenSData[H, Q]{
		MessageID: message.MessageID,
	}

	if message.Header == nil {
		return nil, fmt.Errorf("ravenSParser: header nil")
	} // if

	if output.Header, err = helps.FromProtoAny[H](message.Header); err != nil {
		return nil, fmt.Errorf("ravenSParser: header: %w", err)
	} // if

	if message.Request == nil {
		return nil, fmt.Errorf("ravenSParser: request nil")
	} // if

	if output.Request, err = helps.FromProtoAny[Q](message.Request); err != nil {
		return nil, fmt.Errorf("ravenSParser: request: %w", err)
	} // if

	return output, nil
}

// RavenSData RavenS 資料模型
type RavenSData[H, Q proto.Message] struct {
	MessageID int32 // 訊息編號
	Header    H     // 標頭訊息
	Request   Q     // 請求訊息
}

// Size 取得訊息大小的粗估值
func (this *RavenSData[H, Q]) Size() int {
	size := int(unsafe.Sizeof(this.MessageID))
	size += proto.Size(this.Header)
	size += proto.Size(this.Request)
	return size
}

// Detail 取得詳細資訊字串
func (this *RavenSData[H, Q]) Detail() string {
	builder := &strings.Builder{}
	_, _ = fmt.Fprintf(builder, "messageID: %v\n", this.MessageID)
	_, _ = fmt.Fprintf(builder, "header: %v\n", helps.ProtoString(this.Header))
	_, _ = fmt.Fprintf(builder, "request: %v\n", helps.ProtoString(this.Request))
	return builder.String()
}

// RavenCBuilder 建立 *msgs.RavenC
func RavenCBuilder(messageID, errID int32, header, request proto.Message, respond ...proto.Message) (output any, err error) {
	if header == nil {
		return nil, fmt.Errorf("ravenCBuilder: header nil")
	} // if

	if request == nil {
		return nil, fmt.Errorf("ravenCBuilder: request nil")
	} // if

	message := &msgs.RavenC{
		MessageID: messageID,
		ErrID:     errID,
	}

	if message.Header, err = anypb.New(header); err != nil {
		return nil, fmt.Errorf("ravenCBuilder: header: %w", err)
	} // if

	if message.Request, err = anypb.New(request); err != nil {
		return nil, fmt.Errorf("ravenCBuilder: request: %w", err)
	} // if

	for _, itor := range respond {
		if r, err := anypb.New(itor); err == nil {
			message.Respond = append(message.Respond, r)
		} else {
			return nil, fmt.Errorf("ravenCBuilder: respond: %w", err)
		} // if
	} // for

	return message, nil
}

// RavenCParser 解析 *msgs.RavenC → RavenCData
func RavenCParser[H, Q proto.Message](input any) (output *RavenCData[H, Q], err error) {
	if input == nil {
		return nil, fmt.Errorf("ravenCParser: input nil")
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return nil, fmt.Errorf("ravenCParser: input not *msgs.RavenC")
	} // if

	output = &RavenCData[H, Q]{
		MessageID: message.MessageID,
		ErrID:     message.ErrID,
	}

	if message.Header == nil {
		return nil, fmt.Errorf("ravenCParser: header nil")
	} // if

	if output.Header, err = helps.FromProtoAny[H](message.Header); err != nil {
		return nil, fmt.Errorf("ravenCParser: header: %w", err)
	} // if

	if message.Request == nil {
		return nil, fmt.Errorf("ravenCParser: request nil")
	} // if

	if output.Request, err = helps.FromProtoAny[Q](message.Request); err != nil {
		return nil, fmt.Errorf("ravenCParser: request: %w", err)
	} // if

	for _, itor := range message.Respond {
		if r, err := itor.UnmarshalNew(); err == nil {
			output.Respond = append(output.Respond, r)
		} else {
			return nil, fmt.Errorf("ravenCParser: respond: %w", err)
		} // if
	} // for

	return output, nil
}

// RavenCData RavenC 資料模型
type RavenCData[H, Q proto.Message] struct {
	MessageID int32           // 訊息編號
	ErrID     int32           // 錯誤編號
	Header    H               // 標頭訊息
	Request   Q               // 請求訊息
	Respond   []proto.Message // 回應列表
}

// Size 取得訊息大小的粗估值
func (this *RavenCData[H, Q]) Size() int {
	size := int(unsafe.Sizeof(this.MessageID)) + int(unsafe.Sizeof(this.ErrID))
	size += proto.Size(this.Header)
	size += proto.Size(this.Request)

	for _, itor := range this.Respond {
		size += proto.Size(itor)
	} // for

	return size
}

// Detail 取得詳細資訊字串
func (this *RavenCData[H, Q]) Detail() string {
	builder := &strings.Builder{}
	_, _ = fmt.Fprintf(builder, "messageID: %v\n", this.MessageID)
	_, _ = fmt.Fprintf(builder, "errID: %v\n", this.ErrID)
	_, _ = fmt.Fprintf(builder, "header: %v\n", helps.ProtoString(this.Header))
	_, _ = fmt.Fprintf(builder, "request: %v\n", helps.ProtoString(this.Request))

	for i, itor := range this.Respond {
		_, _ = fmt.Fprintf(builder, "respond[%v]: %v\n", i, helps.ProtoString(itor))
	} // for

	return builder.String()
}

// RavenIsMessageID 檢查 *msgs.RavenC 的 messageID 是否為期望值
func RavenIsMessageID(input any, expected int32) bool {
	message, ok := input.(*msgs.RavenC)
	return ok && message != nil && message.MessageID == expected
}

// RavenIsErrID 檢查 *msgs.RavenC 的 errID 是否為期望值
func RavenIsErrID(input any, expected int32) bool {
	message, ok := input.(*msgs.RavenC)
	return ok && message != nil && message.ErrID == expected
}

// RavenHeader 從 *msgs.RavenC 解析標頭
func RavenHeader[T proto.Message](input any) (result T) {
	if input == nil {
		return result
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return result
	} // if

	result, err := helps.FromProtoAny[T](message.Header)

	if err != nil {
		return result
	} // if

	return result
}

// RavenRequest 從 *msgs.RavenC 解析請求
func RavenRequest[T proto.Message](input any) (result T) {
	if input == nil {
		return result
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return result
	} // if

	result, err := helps.FromProtoAny[T](message.Request)

	if err != nil {
		return result
	} // if

	return result
}

// RavenRespondAt 以索引取得 *msgs.RavenC 回應
func RavenRespondAt[T proto.Message](input any, index int) (result T) {
	if input == nil {
		return result
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return result
	} // if

	size := len(message.Respond)

	if index < 0 || index >= size {
		return result
	} // if

	respond := message.Respond[index]

	if respond == nil {
		return result
	} // if

	result, err := helps.FromProtoAny[T](respond)

	if err != nil {
		return result
	} // if

	return result
}

// RavenRespondFind 取得 *msgs.RavenC 回應列表中第一筆符合型別 T 的回應
func RavenRespondFind[T proto.Message](input any) (result T) {
	if input == nil {
		return result
	} // if

	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return result
	} // if

	messageName := result.ProtoReflect().Descriptor().FullName()

	for _, itor := range message.Respond {
		if itor == nil {
			continue
		} // if

		if itor.MessageName() != messageName {
			continue
		} // if

		if result, err := helps.FromProtoAny[T](itor); err == nil {
			return result
		} // if
	} // for

	return result
}
