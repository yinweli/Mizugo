package msgs

// MJsonQ 要求Json
type MJsonQ struct {
	Time int64 // 傳送時間
}

// MJsonA 回應Json
type MJsonA struct {
	From  *MJsonQ // 來源訊息
	ErrID int     // 錯誤編號
	Count int64   // 封包計數
}

// MLoginQ 要求登入
type MLoginQ struct {
	Account string // 帳號
	Time    int64  // 傳送時間
}

// MLoginA 回應登入
type MLoginA struct {
	From  *MLoginQ // 來源訊息
	ErrID int      // 錯誤編號
	Token string   // 新的token
}

// MUpdateQ 要求更新
type MUpdateQ struct {
	Account string // 帳號
	Token   string // token
	Time    int64  // 傳送時間
}

// MUpdateA 回應更新
type MUpdateA struct {
	From  *MUpdateQ // 來源訊息
	ErrID int       // 錯誤編號
	Token string    // 新的token
}
