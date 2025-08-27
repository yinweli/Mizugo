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

// Err 產生錯誤, 建立一個具「訊息字串 + 錯誤編號」的錯誤物件
//
// 回傳的錯誤型別為 *Error, 其 Error 顯示格式為: <呼叫端函式名>: <訊息串>, ... (<錯誤編號>)
//
// 訊息字串的來源: 依序走訪參數 a, 將可轉成字串的內容(error 或 string)串接起來
// 錯誤編號的來源: 以「最後一個可辨識錯誤編號」為準可辨識來源如下:
//   - *Error: 沿用其內含的錯誤編號
//   - 任何整數/無號整數型別(int/int32/uint64…): 視為錯誤編號
//   - 其他型別將被忽略(不影響訊息與錯誤編號)
//
// 參數建議寫法(方便人讀與機器解析):
//   - Err("說明A", "說明B", err, 錯誤編號)
//   - Err("說明A", "說明B", err) // 無錯誤編號時將使用 ErrUnknown
//   - Err("說明A", 錯誤編號)
//   - Err(err, 錯誤編號)
//   - Err(錯誤編號)
//
// 注意:
//   - 在制定錯誤編號時, 請加入以下三個預設的錯誤編號 Success, ErrUnknown, ErrUnwrap
//   - 為了產出呼叫端位置, Err 會使用 runtime.Caller(1) 取得「呼叫 Err 的函式名稱」, 僅作為訊息前綴
//   - 若 *Error 與一般 error 同時存在, 請確保 *Error 出現在 error 之前(否則可能顯示兩個不同錯誤編號的文字訊息)
//
// 範例:
//
//	if err := doWork(); err != nil {
//	    return Err("處理失敗", err, 10001) // Worker: 處理失敗, <底層訊息> (10001)
//	} // if
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
		case *Error: // *Error 必須在 error 之前, 否則會出現 2 個 errorID 的顯示
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

// UnwrapErrID 從任何 error 物件取出錯誤編號
//
// 若 err 為本套件建立的 *Error, 回傳其錯誤碼; 否則回傳 ErrUnwrap
//
// 範例：
//
//	if err != nil {
//	    switch UnwrapErrID(err) {
//	    case 10001:
//	        // 處理特定錯誤
//
//	    case ErrUnwrap:
//	        // 非 *Error 或無法解碼
//	    } // switch
//	} // if
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
