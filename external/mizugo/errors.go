package mizugo

// 當以mizugo.Errorf產生的錯誤錯誤, 會儲存錯誤編號, 並且此錯誤可以用error傳遞
// 例如:
//     func SomeThing() error {
//         if err := otherFunc(); err != nil {
//             return mizugo.Errorf(1001, err)
//         } // if
//
//         return nil
//     }
// 當需要取得錯誤編號時(例如要回傳給客戶端時), 呼叫mizugo.UnwrapErrorID來取得錯誤
// 當需要取得錯誤訊息時, 仍然可以用error.Error取得, 不過錯誤訊息不會出現錯誤編號
// 例如:
//     errID := UnwrapErrorID(err)

// ErrID 錯誤編號類型
type ErrID = int64

const UnknownErrID ErrID = -1 // 不明的錯誤編號

// Errorf 產生錯誤
func Errorf(errID ErrID, err error) error {
	return &wrapError{
		errID: errID,
		err:   err,
	}
}

// UnwrapErrID 取得錯誤編號, 如果錯誤並非wrapError, 則會獲得UnknownErrID
func UnwrapErrID(err error) ErrID {
	if u, ok := err.(interface {
		ErrID() ErrID
	}); ok {
		return u.ErrID()
	} // if

	return UnknownErrID
}

// wrapError 錯誤結構
type wrapError struct {
	errID ErrID // 錯誤編號
	err   error // 錯誤物件
}

// ErrID 取得錯誤編號
func (this *wrapError) ErrID() ErrID {
	return this.errID
}

// Unwrap 取得錯誤物件
func (this *wrapError) Unwrap() error {
	return this.err
}

// wrapError 取得錯誤訊息
func (this *wrapError) Error() string {
	return this.err.Error()
}
