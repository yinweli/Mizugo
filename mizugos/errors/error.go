package errors

import (
	"fmt"
)

// ErrorID 錯誤編號類型
type ErrorID = int64

const ( // 錯誤編號
	Success ErrorID = iota // 成功
	Unknown                // 不明錯誤
	Max     = 10000        // 最大錯誤編號
)

// Errorf 產生錯誤, 以此方式產生的錯誤內部會儲存錯誤編號, 可以呼叫 UnwrapErrorID 來取得錯誤編號;
// 通常可以在訊息處理時使用此錯誤工具, 外部系統如果要編製錯誤編號, 需從 Max 之後開始編制
func Errorf(errorID ErrorID, err error) error {
	return &wrapError{
		errorID: errorID,
		err:     err,
	}
}

// UnwrapErrorID 取得錯誤編號, 如果錯誤並非由 Errorf 產生的, 則會獲得 Unknown
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
