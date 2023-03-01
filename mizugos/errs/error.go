package errs

import (
	"fmt"
)

// Errort 產生錯誤, 以此方式產生的錯誤內容會包含呼叫字串, 自訂標籤
func Errort(tag any) error {
	return &wrapError{
		tag: tag,
		err: nil,
	}
}

// Errore 產生錯誤, 以此方式產生的錯誤內容會包含呼叫字串, 自訂標籤, 錯誤內容
func Errore(tag any, err error) error {
	return &wrapError{
		tag: tag,
		err: err,
	}
}

// Errorf 產生錯誤, 以此方式產生的錯誤內容會包含呼叫字串, 自訂標籤, 錯誤內容
func Errorf(tag any, format string, a ...any) error {
	return &wrapError{
		tag: tag,
		err: fmt.Errorf(format, a...),
	}
}

// wrapError 錯誤結構
type wrapError struct {
	tag any   // 標籤資料
	err error // 錯誤物件
}

// Error 取得錯誤訊息
func (this *wrapError) Error() string {
	if this.err == nil {
		return fmt.Sprintf("[%v]", this.tag)
	} else {
		return fmt.Sprintf("[%v] %v", this.tag, this.err.Error())
	} // if
}
