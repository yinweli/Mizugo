package helpers

import (
	"fmt"
)

// 當以helpers.Errorf產生的錯誤錯誤, 會儲存錯誤編號, 並且此錯誤可以用error傳遞
//     func something() error {
//         if err := otherFunc(); err != nil {
//             return helpers.Errorf(1001, err)
//         } // if
//
//         return nil
//     }
// 當需要取得錯誤編號時(例如要回傳給客戶端時), 呼叫helpers.UnwrapErrorID來取得錯誤編號
//     errorID := UnwrapErrorID(err)
// 當需要取得錯誤訊息時(例如要輸出日誌時), 可以用err.Error()取得
//     err := helpers.Errorf(1001, otherError)
//     message := err.Error()

// ErrorID 錯誤編號類型
type ErrorID = int64

const (
	Unknown ErrorID = -1 // 不明錯誤
	Success ErrorID = 0  // 成功
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
	return fmt.Sprintf("[%d] %s", this.errorID, this.err.Error())
}
