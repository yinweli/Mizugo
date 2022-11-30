package logs

// empty提供空日誌, 以此方式輸出的日誌都會被丟棄, 也不會顯示在控制台上, 這是預設的日誌模式
// 使用方式:
//   執行logs.Set(logs.NewEmpty)

// NewEmpty 建立空日誌
func NewEmpty(_ string, _ Level) Logger {
	return &Empty{}
}

// Empty 空日誌
type Empty struct {
}

// Message 記錄訊息
func (this *Empty) Message(_ string) Logger {
	return this
}

// KV 記錄索引與數值
func (this *Empty) KV(_ string, _ any) Logger {
	return this
}

// Error 記錄錯誤
func (this *Empty) Error(_ error) Logger {
	return this
}

// EndError 以錯誤結束記錄
func (this *Empty) EndError(err error) error {
	return err
}

// End 結束記錄
func (this *Empty) End() {
}
