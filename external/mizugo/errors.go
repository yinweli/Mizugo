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
// 當需要取出錯誤編號時(例如要回傳給客戶端時), 呼叫mizugo.UnwrapErrorID來取得錯誤
// 例如:
//     errID := UnwrapErrorID(err)

// Errorf 產生錯誤
func Errorf(errID int64, err error) error {
    return &wrapError{
        errID: errID,
        err:   err,
    }
}

// UnwrapErrID 取出錯誤編號, 如果錯誤並非wrapError, 則會獲得0
func UnwrapErrID(err error) int64 {
    if u, ok := err.(interface {
        ErrID() int64
    }); ok {
        return u.ErrID()
    } // if

    return 0
}

// wrapError 錯誤結構
type wrapError struct {
    errID int64 // 錯誤編號
    err   error // 錯誤物件
}

// ErrID 取得錯誤編號
func (this *wrapError) ErrID() int64 {
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
