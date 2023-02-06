package logs

// 空日誌, 作為預設的日誌使用, 不會輸出任何訊息

// EmptyLogger 空日誌
type EmptyLogger struct {
}

// Initialize 初始化處理
func (this *EmptyLogger) Initialize() error {
	return nil
}

// Finalize 結束處理
func (this *EmptyLogger) Finalize() {
}

// New 建立日誌
func (this *EmptyLogger) New(_ string, _ Level) Stream {
	return &EmptyStream{}
}

// EmptyStream 空記錄
type EmptyStream struct {
}

// Message 記錄訊息
func (this *EmptyStream) Message(_ string, _ ...any) Stream {
	return this
}

// KV 記錄索引與數值
func (this *EmptyStream) KV(_ string, _ any) Stream {
	return this
}

// Error 記錄錯誤
func (this *EmptyStream) Error(_ error) Stream {
	return this
}

// EndError 以錯誤結束記錄
func (this *EmptyStream) EndError(_ error) {
}

// End 結束記錄
func (this *EmptyStream) End() {
}