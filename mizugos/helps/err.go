package helps

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

const (
	Success    = 0 // 成功
	ErrUnknown = 1 // 不明錯誤
	ErrUnwrap  = 2 // 取得錯誤編號失敗
)

// Err 產生錯誤, 最後會產生 Error 物件, 物件中分為字串與錯誤編號兩個部分, 產生的方式為
//   - 字串部分: 由 a 列表中的項目轉為字串後組合而成
//   - 錯誤編號部分: a 列表中最後一個可獲取的錯誤編號, 只有項目為錯誤編號或是 Error 才能獲取錯誤編號
//
// 使用時, 請讓參數結尾為 [ error | 字串 | Error, 錯誤編號]; 以下提供使用的範例
//   - trials.Err(字串, 字串, ... , error, 錯誤編號)
//   - trials.Err(字串, 字串, ... , error)
//   - trials.Err(字串, 字串, ... , 錯誤編號)
//   - trials.Err(字串, 錯誤編號)
//   - trials.Err(error, 錯誤編號)
//   - trials.Err(錯誤編號)
func Err(a ...any) error {
	builder := strings.Builder{}
	errorID := ErrUnknown
	separate := false

	if pc, _, _, ok := runtime.Caller(1); ok {
		builder.WriteString(filepath.Base(runtime.FuncForPC(pc).Name()))
	} else {
		builder.WriteString("unknown")
	} // if

	builder.WriteString(": ")

	for _, itor := range a {
		sz := ""

		switch v := itor.(type) {
		case *Error: // *Error 必須在 error 之前, 否則會出現2個 errorID 的顯示
			sz = v.err
			errorID = v.errID

		case error:
			sz = v.Error()

		case string:
			sz = v

		default:
			errorID = convertToInt(reflect.ValueOf(v))
		} // switch

		if sz == "" {
			continue
		} // if

		if separate {
			builder.WriteString(", ")
		} else {
			separate = true
		} // if

		builder.WriteString(sz)
	} // for

	return &Error{
		err:   builder.String(),
		errID: errorID,
	}
}

// UnwrapErrID 從錯誤物件取得錯誤編號
func UnwrapErrID(err error) int {
	var e *Error

	if errors.As(err, &e) {
		return e.errID
	} // if

	return ErrUnwrap
}

// Error 錯誤資料
type Error struct {
	err   string // 錯誤字串
	errID int    // 錯誤編號
}

// Error 取得錯誤字串
func (this Error) Error() string {
	return fmt.Sprintf("%v (%v)", this.err, this.errID)
}

// convertToInt 轉換物件為int
func convertToInt(v reflect.Value) int {
	switch v.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(v.Uint())

	default:
		return ErrUnknown
	} // switch
}
