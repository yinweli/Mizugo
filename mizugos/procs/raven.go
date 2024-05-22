package procs

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/msgs"
)

// NewRaven 建立raven處理器
func NewRaven() *Raven {
	return &Raven{
		Procmgr: NewProcmgr(),
	}
}

// Raven raven處理器, 封包結構使用msgs.RavenS, msgs.RavenC
//   - 訊息定義: support/proto/mizugo/raven.proto
type Raven struct {
	*Procmgr // 管理器
}

// Encode 封包編碼
func (this *Raven) Encode(input any) (output any, err error) {
	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return nil, fmt.Errorf("raven encode: input")
	} // if

	output, err = proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("raven encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Raven) Decode(input any) (output any, err error) {
	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("raven decode: input")
	} // if

	message := &msgs.RavenS{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("raven decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Raven) Process(input any) error {
	message, ok := input.(*msgs.RavenS)

	if ok == false {
		return fmt.Errorf("raven process: input")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("raven process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// NewRavenClient 建立raven客戶端處理器
func NewRavenClient() *RavenClient {
	return &RavenClient{
		Procmgr: NewProcmgr(),
	}
}

// RavenClient raven客戶端處理器, 封包結構使用msgs.RavenS, msgs.RavenC; 這個處理器提供給客戶端使用
//   - 訊息定義: support/proto/mizugo/raven.proto
type RavenClient struct {
	*Procmgr // 管理器
}

// Encode 封包編碼
func (this *RavenClient) Encode(input any) (output any, err error) {
	message, ok := input.(*msgs.RavenS)

	if ok == false {
		return nil, fmt.Errorf("raven client encode: input")
	} // if

	output, err = proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("raven client encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *RavenClient) Decode(input any) (output any, err error) {
	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("raven client decode: input")
	} // if

	message := &msgs.RavenC{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("raven client decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *RavenClient) Process(input any) error {
	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return fmt.Errorf("raven client process: input")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("raven client process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// RavenSBuilder 建立RavenS訊息
func RavenSBuilder(messageID int32, header, request proto.Message) (output any, err error) {
	message := &msgs.RavenS{}
	message.MessageID = messageID

	if message.Header, err = anypb.New(header); err != nil {
		return nil, fmt.Errorf("ravenSBuilder: header: %w", err)
	} // if

	if message.Request, err = anypb.New(request); err != nil {
		return nil, fmt.Errorf("ravenSBuilder: request: %w", err)
	} // if

	return message, nil
}

// RavenSParser 解析RavenS訊息
func RavenSParser[H, Q any](input any) (output *RavenSData[H, Q], err error) {
	message, ok := input.(*msgs.RavenS)

	if ok == false {
		return nil, fmt.Errorf("ravenSParser: input")
	} // if

	output = &RavenSData[H, Q]{}
	output.MessageID = message.MessageID

	if output.Header, err = helps.FromProtoAny[H](message.Header); err != nil {
		return nil, fmt.Errorf("ravenSParser: header: %w", err)
	} // if

	if output.Request, err = helps.FromProtoAny[Q](message.Request); err != nil {
		return nil, fmt.Errorf("ravenSParser: request: %w", err)
	} // if

	return output, nil
}

// RavenCBuilder 建立RavenC訊息
func RavenCBuilder(messageID, errID int32, header, request proto.Message, respond ...proto.Message) (output any, err error) {
	message := &msgs.RavenC{}
	message.MessageID = messageID
	message.ErrID = errID

	if message.Header, err = anypb.New(header); err != nil {
		return nil, fmt.Errorf("ravenCBuilder: header: %w", err)
	} // if

	if message.Request, err = anypb.New(request); err != nil {
		return nil, fmt.Errorf("ravenCBuilder: request: %w", err)
	} // if

	for _, itor := range respond {
		if temp, err := anypb.New(itor); err == nil {
			message.Respond = append(message.Respond, temp)
		} else {
			return nil, fmt.Errorf("ravenCBuilder: respond: %w", err)
		} // if
	} // for

	return message, nil
}

// RavenCParser 解析RavenC訊息
func RavenCParser[H, Q any](input any) (output *RavenCData[H, Q], err error) {
	message, ok := input.(*msgs.RavenC)

	if ok == false {
		return nil, fmt.Errorf("ravenCParser: input")
	} // if

	output = &RavenCData[H, Q]{}
	output.MessageID = message.MessageID
	output.ErrID = message.ErrID

	if output.Header, err = helps.FromProtoAny[H](message.Header); err != nil {
		return nil, fmt.Errorf("ravenCParser: header: %w", err)
	} // if

	if output.Request, err = helps.FromProtoAny[Q](message.Request); err != nil {
		return nil, fmt.Errorf("ravenCParser: request: %w", err)
	} // if

	for _, itor := range message.Respond {
		if temp, err := itor.UnmarshalNew(); err == nil {
			output.Respond = append(output.Respond, temp)
		} else {
			return nil, fmt.Errorf("ravenCParser: respond: %w", err)
		} // if
	} // for

	return output, nil
}

// RavenSData RavenS資料
type RavenSData[H, Q any] struct {
	MessageID int32 // 訊息編號
	Header    *H    // 標頭資料
	Request   *Q    // 要求資料
}

// Size 取得訊息大小
func (this *RavenSData[H, Q]) Size() int {
	size := int(unsafe.Sizeof(this.MessageID))
	size += proto.Size(any(this.Header).(proto.Message))
	size += proto.Size(any(this.Request).(proto.Message))
	return size
}

// Detail 取得訊息資訊
func (this *RavenSData[H, Q]) Detail() string {
	builder := &strings.Builder{}
	_, _ = fmt.Fprintf(builder, "messageID: %v\n", this.MessageID)
	_, _ = fmt.Fprintf(builder, "header: %v\n", helps.ProtoString(any(this.Header).(proto.Message)))
	_, _ = fmt.Fprintf(builder, "request: %v\n", helps.ProtoString(any(this.Request).(proto.Message)))
	return builder.String()
}

// RavenCData RavenC資料
type RavenCData[H, Q any] struct {
	MessageID int32           // 訊息編號
	ErrID     int32           // 錯誤編號
	Header    *H              // 標頭資料
	Request   *Q              // 要求資料
	Respond   []proto.Message // 回應列表
}

// Size 取得訊息大小
func (this *RavenCData[H, Q]) Size() int {
	size := int(unsafe.Sizeof(this.MessageID)) + int(unsafe.Sizeof(this.ErrID))
	size += proto.Size(any(this.Header).(proto.Message))
	size += proto.Size(any(this.Request).(proto.Message))

	for _, itor := range this.Respond {
		size += proto.Size(any(itor).(proto.Message))
	} // for

	return size
}

// Detail 取得詳細資訊
func (this *RavenCData[H, Q]) Detail() string {
	builder := &strings.Builder{}
	_, _ = fmt.Fprintf(builder, "messageID: %v\n", this.MessageID)
	_, _ = fmt.Fprintf(builder, "errID: %v\n", this.ErrID)
	_, _ = fmt.Fprintf(builder, "header: %v\n", helps.ProtoString(any(this.Header).(proto.Message)))
	_, _ = fmt.Fprintf(builder, "request: %v\n", helps.ProtoString(any(this.Request).(proto.Message)))

	for i, itor := range this.Respond {
		_, _ = fmt.Fprintf(builder, "respond[%v]: %v\n", i, helps.ProtoString(itor))
	} // for

	return builder.String()
}

// RavenRespond 取得回應列表中首個符合指定類型的資料
func RavenRespond[T any](respond []proto.Message) *T {
	for _, itor := range respond {
		if obj, ok := any(itor).(*T); ok {
			return obj
		} // if
	} // for

	return nil
}

// RavenRespondAt 取得回應列表中指定位置的資料
func RavenRespondAt[T any](respond []proto.Message, index int) *T {
	if index < len(respond) {
		if obj, ok := any(respond[index]).(*T); ok {
			return obj
		} // if
	} // if

	return nil
}

// RavenTestMessageID 測試訊息編號是否相符, input必須是msgs.RavenC, 並且訊息編號與expected相符才會傳回true, 否則為false
func RavenTestMessageID(input any, expected int32) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("messageID: invalid input\n")
		return false
	} // if

	if expected != actual.MessageID {
		fmt.Printf("messageID: expected {%v} but actual {%v}\n", expected, actual.MessageID)
		return false
	} // if

	return true
}

// RavenTestErrID 測試錯誤編號是否相符, input必須是msgs.RavenC, 並且錯誤編號與expected相符才會傳回true, 否則為false
func RavenTestErrID(input any, expected int32) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("errID: invalid input\n")
		return false
	} // if

	if expected != actual.ErrID {
		fmt.Printf("errID: expected {%v} but actual {%v}\n", expected, actual.ErrID)
		return false
	} // if

	return true
}

// RavenTestHeader 測試標頭資料是否相符, input必須是msgs.RavenC, 並且標頭資料與expected相符才會傳回true, 否則為false
func RavenTestHeader(input any, expected proto.Message) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("header: invalid input\n")
		return false
	} // if

	header, err := actual.Header.UnmarshalNew()

	if err != nil {
		fmt.Printf("header: unmarshal failed\n")
		return false
	} // if

	if proto.Equal(expected, header) == false {
		fmt.Printf("header: expected {%v} but actual {%v}\n", expected, header)
		return false
	} // if

	return true
}

// RavenTestRequest 測試要求資料是否相符, input必須是msgs.RavenC, 並且要求資料與expected相符才會傳回true, 否則為false
func RavenTestRequest(input any, expected proto.Message) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("request: invalid input\n")
		return false
	} // if

	request, err := actual.Request.UnmarshalNew()

	if err != nil {
		fmt.Printf("request: unmarshal failed\n")
		return false
	} // if

	if proto.Equal(expected, request) == false {
		fmt.Printf("request: expected {%v} but actual {%v}\n", expected, request)
		return false
	} // if

	return true
}

// RavenTestRespond 測試回應列表資料是否相符, input必須是msgs.RavenC, 並且expected列表的每一項元素都有與之相符的回應資料才會傳回true, 否則為false;
// 但是這並不代表expected列表與回應列表完全一致, 例如有只出現於回應列表, 但不在expected列表中的資料, 就無法通過此函式檢測出來
func RavenTestRespond(input any, expected ...proto.Message) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("respond: invalid input\n")
		return false
	} // if

	result := true
	report := &strings.Builder{}
	report.WriteString("respond:\n")

	for i, itor := range expected {
		if i >= len(actual.Respond) {
			result = false
			_, _ = fmt.Fprintf(report, "    %v not found\n", i)
			continue
		} // if

		respond, err := actual.Respond[i].UnmarshalNew()

		if err != nil {
			result = false
			_, _ = fmt.Fprintf(report, "    %v unmarshal failed\n", i)
			continue
		} // if

		if proto.Equal(itor, respond) == false {
			result = false
			_, _ = fmt.Fprintf(report, "    %v expected {%v} but actual {%v}\n", i, itor, respond)
			continue
		} // if
	} // for

	if result == false {
		fmt.Print(report.String())
		return false
	} // if

	return true
}

// RavenTestRespondType 測試回應列表類型是否相符, input必須是msgs.RavenC, 並且expected列表的每一項元素都有與之相符的回應類型才會傳回true, 否則為false;
// 但是這並不代表expected列表與回應列表完全一致, 例如有只出現於回應列表, 但不在expected列表中的類型, 就無法通過此函式檢測出來
func RavenTestRespondType(input any, expected ...proto.Message) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("respond type: invalid input\n")
		return false
	} // if

	result := true
	report := &strings.Builder{}
	report.WriteString("respond type:\n")

	for i, itor := range expected {
		if i >= len(actual.Respond) {
			result = false
			_, _ = fmt.Fprintf(report, "    %v not found\n", i)
			continue
		} // if

		respond, err := actual.Respond[i].UnmarshalNew()

		if err != nil {
			result = false
			_, _ = fmt.Fprintf(report, "    %v unmarshal failed\n", i)
			continue
		} // if

		expectedType := reflect.TypeOf(itor).String()
		actualType := reflect.TypeOf(respond).String()

		if expectedType != actualType {
			result = false
			_, _ = fmt.Fprintf(report, "    %v expected {%v} but actual {%v}\n", i, expectedType, actualType)
			continue
		} // if
	} // for

	if result == false {
		fmt.Print(report.String())
		return false
	} // if

	return true
}

// RavenTestRespondLength 測試回應列表長度是否相符, input必須是msgs.RavenC, 並且回應列表長度必須相符才會傳回true, 否則為false
func RavenTestRespondLength(input any, expected int) bool {
	actual, ok := input.(*msgs.RavenC)

	if ok == false {
		fmt.Printf("length: invalid input\n")
		return false
	} // if

	if length := len(actual.Respond); expected != length {
		fmt.Printf("length: expected {%v} but actual {%v}\n", expected, length)
		return false
	} // if

	return true
}
