package msgs

// NewStringMsg 建立字串訊息器
func NewStringMsg() *StringMsg {
	return &StringMsg{
		Msgmgr: NewMsgmgr(),
	}
}

// StringMsg 字串訊息器
type StringMsg struct {
	*Msgmgr // 訊息管理器
}

// Encode 封包編碼
func (this *StringMsg) Encode(message any) (packet []byte, err error) {
	// TODO: 檢查是否StringMessage, md5, 轉json字串, 轉[]byte, 加密
	return nil, nil
}

// Decode 封包解碼
func (this *StringMsg) Decode(packet []byte) (message any, err error) {
	// TODO: 解密, 轉json(StringMessage), md5驗證
	return nil, nil
}

// Process 訊息處理
func (this *StringMsg) Process(message any) error {
	// TODO: 檢查是否StringMessage, 用訊息編號取得處理函式, 執行處理!!!
	return nil
}

// StringMessage 字串訊息資料
type StringMessage struct {
	MessageID MessageID `json:"messageID"` // 訊息編號
	Message   string    `json:"message"`   // 訊息字串
	MD5       string    `json:"MD5"`       // 訊息驗證
}
