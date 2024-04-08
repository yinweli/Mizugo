package helps

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ExpvarStr 取得度量字串
func ExpvarStr(expvarStat []ExpvarStat) string {
	builder := &strings.Builder{}
	builder.WriteByte('{')
	first := true

	for _, itor := range expvarStat {
		if first == false {
			builder.WriteString(", ")
		} // if

		if itor.stringType() {
			_, _ = fmt.Fprintf(builder, "\"%v\": \"%v\"", itor.Name, itor.Data)
		} else {
			_, _ = fmt.Fprintf(builder, "\"%v\": %v", itor.Name, itor.Data)
		} // if

		first = false
	} // for

	builder.WriteByte('}')
	return builder.String()
}

// ExpvarStat 度量資料
type ExpvarStat struct {
	Name string // 名稱
	Data any    // 資料
}

// stringValue 取得資料是否為字串類型
func (this ExpvarStat) stringType() bool {
	switch this.Data.(type) {
	case nil, string, time.Duration:
		return true

	default:
		return reflect.TypeOf(this.Data).Kind() == reflect.Struct
	} // switch
}
