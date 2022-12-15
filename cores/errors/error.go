package errors

import (
	"fmt"
)

// 當以errors.Errorf產生的錯誤, 會儲存錯誤編號, 並且此錯誤可以用error傳遞
//     if err := otherFunc(); err != nil {
//         return errors.Errorf(1001, err)
//     } // if
// 當需要取得錯誤編號時(例如要回傳給客戶端時), 呼叫errors.UnwrapErrorID來取得錯誤編號
//     errorID := UnwrapErrorID(err)
// 當需要取得錯誤訊息時(例如要輸出日誌時), 可以用err.Error()取得
//     err := errors.Errorf(1001, otherError)
//     message := err.Error()
// 通常會在封包接收層使用此錯誤工具

// ErrorID 錯誤編號類型
type ErrorID = int64

// 錯誤編號列表
const (
	Success ErrorID  = iota // 成功
	Unknown                 // 不明錯誤
	Max     = 100000        // 最大錯誤編號, 外部系統由此編號之後編制自己的錯誤編號
)

// Errorf 產生錯誤
func Errorf(errorID ErrorID, err error) error {
	return &wrapError{
		errorID: errorID,
		err:     err,
	}
}

// UnwrapErrorID 取得錯誤編號, 如果錯誤並非wrapError, 則會獲得Unknown
func UnwrapErrorID(err error) ErrorID {
	if u, ok := err.(interface {
		ErrorID() ErrorID
	}); ok {
		return u.ErrorID()
	} // if

	return Unknown
}

// wrapError 錯誤結構
type wrapError struct {
	errorID ErrorID // 錯誤編號
	err     error   // 錯誤物件
}

// ErrorID 取得錯誤編號
func (this *wrapError) ErrorID() ErrorID {
	return this.errorID
}

// Unwrap 取得錯誤物件
func (this *wrapError) Unwrap() error {
	return this.err
}

// wrapError 取得錯誤訊息
func (this *wrapError) Error() string {
	return fmt.Sprintf("[%v] %v", this.errorID, this.err.Error())
}
