package patterns

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/procs"
)

// Raven系列工具專為伺服器與客戶端之間的訊息傳遞協議設計
// 它基於 procs/proto.go 中的 proto 處理器, 利用 msgs.RavenQData 和 msgs.RavenAData 對基礎協議的訊息內容進行細化處理
// 其中, 客戶端向伺服器發送的訊息採用 msgs.RavenQData 格式, 而伺服器向客戶端返回的訊息則採用 msgs.RavenAData 格式

// RavenQBuilder 建立RavenQ訊息
func RavenQBuilder(messageID procs.MessageID, header, request proto.Message) (output any, err error) {
	message := &msgs.RavenQ{}

	if message.Header, err = anypb.New(header); err != nil {
		return nil, fmt.Errorf("ravenQBuilder: header: %w", err)
	} // if

	if message.Request, err = anypb.New(request); err != nil {
		return nil, fmt.Errorf("ravenQBuilder: request: %w", err)
	} // if

	if output, err = procs.ProtoMarshal(messageID, message); err != nil {
		return nil, fmt.Errorf("ravenQBuilder: %w", err)
	} // if

	return output, nil
}

// RavenQParser 解析RavenQ訊息
func RavenQParser[H, Q any](input any) (output *RavenQData[H, Q], err error) {
	output = &RavenQData[H, Q]{}
	message := (*msgs.RavenQ)(nil)

	if output.MessageID, message, err = procs.ProtoUnmarshal[msgs.RavenQ](input); err != nil {
		return nil, fmt.Errorf("ravenQParser: %w", err)
	} // if

	if output.Header, err = procs.ProtoAny[H](message.Header); err != nil {
		return nil, fmt.Errorf("ravenQParser: header: %w", err)
	} // if

	if output.Request, err = procs.ProtoAny[Q](message.Request); err != nil {
		return nil, fmt.Errorf("ravenQParser: output: %w", err)
	} // if

	return output, nil
}

// RavenABuilder 建立RavenA訊息
func RavenABuilder(messageID procs.MessageID, errID int32, header, request proto.Message, respond ...proto.Message) (output any, err error) {
	message := &msgs.RavenA{}
	message.ErrID = errID

	if message.Header, err = anypb.New(header); err != nil {
		return nil, fmt.Errorf("ravenABuilder: header: %w", err)
	} // if

	if message.Request, err = anypb.New(request); err != nil {
		return nil, fmt.Errorf("ravenABuilder: request: %w", err)
	} // if

	for _, itor := range respond {
		if temp, err := anypb.New(itor); err == nil {
			message.Respond = append(message.Respond, temp)
		} else {
			return nil, fmt.Errorf("ravenABuilder: respond: %w", err)
		} // if
	} // for

	if output, err = procs.ProtoMarshal(messageID, message); err != nil {
		return nil, fmt.Errorf("ravenABuilder: %w", err)
	} // if

	return output, nil
}

// RavenAParser 解析RavenA訊息
func RavenAParser[H, Q any](input any) (output *RavenAData[H, Q], err error) {
	output = &RavenAData[H, Q]{}
	message := (*msgs.RavenA)(nil)

	if output.MessageID, message, err = procs.ProtoUnmarshal[msgs.RavenA](input); err != nil {
		return nil, fmt.Errorf("ravenAParser: %w", err)
	} // if

	output.ErrID = message.ErrID

	if output.Header, err = procs.ProtoAny[H](message.Header); err != nil {
		return nil, fmt.Errorf("ravenAParser: header: %w", err)
	} // if

	if output.Request, err = procs.ProtoAny[Q](message.Request); err != nil {
		return nil, fmt.Errorf("ravenAParser: request: %w", err)
	} // if

	for _, itor := range message.Respond {
		if temp, err := itor.UnmarshalNew(); err == nil {
			output.Respond = append(output.Respond, temp)
		} else {
			return nil, fmt.Errorf("ravenAParser: respond: %w", err)
		} // if
	} // for

	return output, nil
}

// RavenTestMessageID 測試訊息編號是否相符;
// input必須是msgs.ProtoMsg並且訊息編號與expected相符才會傳回true, 否則為false
func RavenTestMessageID(input any, expected procs.MessageID) bool {
	if actual, ok := input.(*msgs.ProtoMsg); ok {
		return expected == actual.MessageID
	} // if

	return false
}

// RavenTestErrID 測試錯誤編號是否相符;
// input必須是msgs.ProtoMsg並且內存msgs.RavenAData, 並且錯誤編號與expected相符才會傳回true, 否則為false
func RavenTestErrID(input any, expected int32) bool {
	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenA](input); err == nil {
		return expected == actual.ErrID
	} // if

	return false
}

// RavenTestHeader 測試標頭資料是否相符;
// input必須是msgs.ProtoMsg並且內存msgs.RavenQ或是msgs.RavenAData, 並且標頭資料與expected相符才會傳回true, 否則為false
func RavenTestHeader(input any, expected proto.Message) bool {
	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenQ](input); err == nil {
		if header, err := actual.Header.UnmarshalNew(); err == nil {
			return proto.Equal(expected, header)
		} // if
	} // if

	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenA](input); err == nil {
		if header, err := actual.Header.UnmarshalNew(); err == nil {
			return proto.Equal(expected, header)
		} // if
	} // if

	return false
}

// RavenTestRequest 測試要求資料是否相符;
// input必須是msgs.ProtoMsg並且內存msgs.RavenQ或是msgs.RavenAData, 並且要求資料與expected相符才會傳回true, 否則為false
func RavenTestRequest(input any, expected proto.Message) bool {
	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenQ](input); err == nil {
		if request, err := actual.Request.UnmarshalNew(); err == nil {
			return proto.Equal(expected, request)
		} // if
	} // if

	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenA](input); err == nil {
		if request, err := actual.Request.UnmarshalNew(); err == nil {
			return proto.Equal(expected, request)
		} // if
	} // if

	return false
}

// RavenTestRespond 測試回應列表資料是否相符;
// input必須是msgs.ProtoMsg並且內存msgs.RavenAData, 並且expected列表的每一項元素都有與之相符的回應資料才會傳回true, 否則為false;
// 但是這並不代表expected列表與回應列表完全一致, 例如有只出現於回應列表, 但不在expected列表中的資料, 就無法通過此函式檢測出來
func RavenTestRespond(input any, expected ...proto.Message) bool {
	fmt.Printf(">> test respond\n")

	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenA](input); err == nil {
		result := true

		for i, itor := range expected {
			expectedType := reflect.TypeOf(itor).String()
			actualType := "unknown" //nolint:goconst
			match := false

			if i < len(actual.Respond) {
				if respond, err := actual.Respond[i].UnmarshalNew(); err == nil {
					actualType = reflect.TypeOf(respond).String()
					match = proto.Equal(itor, respond)
				} // if
			} // if

			result = result && match
			fmt.Printf("    [%4v] %v %v %v\n", i, expectedType, bool2symbol(match), actualType)
		} // for

		return result
	} // if

	return false
}

// RavenTestRespondType 測試回應列表類型是否相符;
// input必須是msgs.ProtoMsg並且內存msgs.RavenAData, 並且expected列表的每一項元素都有與之相符的回應類型才會傳回true, 否則為false;
// 但是這並不代表expected列表與回應列表完全一致, 例如有只出現於回應列表, 但不在expected列表中的類型, 就無法通過此函式檢測出來
func RavenTestRespondType(input any, expected ...proto.Message) bool {
	fmt.Printf(">> test respond type\n")

	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenA](input); err == nil {
		result := true

		for i, itor := range expected {
			expectedType := reflect.TypeOf(itor).String()
			actualType := "unknown"
			match := false

			if i < len(actual.Respond) {
				if respond, err := actual.Respond[i].UnmarshalNew(); err == nil {
					actualType = reflect.TypeOf(respond).String()
					match = expectedType == actualType
				} // if
			} // if

			result = result && match
			fmt.Printf("    [%4v] %v %v %v\n", i, expectedType, bool2symbol(match), actualType)
		} // for

		return result
	} // if

	return false
}

// RavenTestRespondLength 測試回應列表長度是否相符;
// input必須是msgs.ProtoMsg並且內存msgs.RavenAData, 並且回應列表長度必須相符才會傳回true, 否則為false
func RavenTestRespondLength(input any, expected int) bool {
	if _, actual, err := procs.ProtoUnmarshal[msgs.RavenA](input); err == nil {
		return expected == len(actual.Respond)
	} // if

	return false
}

// RavenQData RavenQ資料
type RavenQData[H, Q any] struct {
	MessageID procs.MessageID // 訊息編號
	Header    *H              // 標頭資料
	Request   *Q              // 要求資料
}

// Detail 取得詳細資訊
func (this *RavenQData[H, Q]) Detail() string {
	leno := int(unsafe.Sizeof(this.MessageID))
	lenh := proto.Size(any(this.Header).(proto.Message))
	lenq := proto.Size(any(this.Request).(proto.Message))
	builder := &strings.Builder{}
	_, _ = fmt.Fprintf(builder, ">> message\n")
	_, _ = fmt.Fprintf(builder, "    messageID: %v\n", this.MessageID)
	_, _ = fmt.Fprintf(builder, "    header: %v\n", proto2json(this.Header))
	_, _ = fmt.Fprintf(builder, "    request: %v\n", proto2json(this.Request))
	_, _ = fmt.Fprintf(builder, ">> size\n")
	_, _ = fmt.Fprintf(builder, "    other: %v\n", leno)
	_, _ = fmt.Fprintf(builder, "    header: %v\n", lenh)
	_, _ = fmt.Fprintf(builder, "    request: %v\n", lenq)
	_, _ = fmt.Fprintf(builder, "    total: %v\n", leno+lenh+lenq)
	return builder.String()
}

// RavenAData RavenA資料
type RavenAData[H, Q any] struct {
	MessageID procs.MessageID // 訊息編號
	ErrID     int32           // 錯誤編號
	Header    *H              // 標頭資料
	Request   *Q              // 要求資料
	Respond   []proto.Message // 回應列表
}

// GetRespond 取得回應列表中首個符合指定類型的資料
func (this *RavenAData[H, Q]) GetRespond(type_ proto.Message) proto.Message {
	if type_ != nil {
		for _, itor := range this.Respond {
			if itor.ProtoReflect().Type() == type_.ProtoReflect().Type() {
				return itor
			} // if
		} // for
	} // if

	return nil
}

// GetRespondAt 取得回應列表中指定位置的資料
func (this *RavenAData[H, Q]) GetRespondAt(index int) proto.Message {
	if index < len(this.Respond) {
		return this.Respond[index]
	} // if

	return nil
}

// Detail 取得詳細資訊
func (this *RavenAData[H, Q]) Detail() string {
	leno := int(unsafe.Sizeof(this.MessageID)) + int(unsafe.Sizeof(this.ErrID))
	lenh := proto.Size(any(this.Header).(proto.Message))
	lenq := proto.Size(any(this.Request).(proto.Message))
	lent := leno + lenh + lenq
	builder := &strings.Builder{}
	_, _ = fmt.Fprintf(builder, ">> message\n")
	_, _ = fmt.Fprintf(builder, "    messageID: %v\n", this.MessageID)
	_, _ = fmt.Fprintf(builder, "    errID: %v\n", this.ErrID)
	_, _ = fmt.Fprintf(builder, "    header: %v\n", proto2json(this.Header))
	_, _ = fmt.Fprintf(builder, "    request: %v\n", proto2json(this.Request))

	for i, itor := range this.Respond {
		_, _ = fmt.Fprintf(builder, "    respond[%v]: %v\n", i, proto2json(itor))
	} // for

	_, _ = fmt.Fprintf(builder, ">> size\n")
	_, _ = fmt.Fprintf(builder, "    other: %v\n", leno)
	_, _ = fmt.Fprintf(builder, "    header: %v\n", lenh)
	_, _ = fmt.Fprintf(builder, "    request: %v\n", lenq)

	for i, itor := range this.Respond {
		lens := proto.Size(any(itor).(proto.Message))
		lent += lens
		_, _ = fmt.Fprintf(builder, "    respond[%v]: %v\n", i, lens)
	} // for

	_, _ = fmt.Fprintf(builder, "    total: %v\n", lent)
	return builder.String()
}

// proto2json 將proto轉為json字串
func proto2json(input any) string {
	json, _ := protojson.Marshal(input.(proto.Message))
	return string(json)
}

// bool2symbol 將bool轉為比較字串
func bool2symbol(input bool) string {
	if input {
		return "=="
	} else {
		return "!="
	} // if
}
