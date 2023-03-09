package logs

// EmptyLogger 空日誌, 作為預設的日誌使用, 不會輸出任何訊息
type EmptyLogger struct {
}

// Initialize 初始化處理
func (this *EmptyLogger) Initialize() error {
	return nil
}

// Finalize 結束處理
func (this *EmptyLogger) Finalize() {
}

// Debug 記錄除錯訊息, 用於記錄除錯訊息
func (this *EmptyLogger) Debug(_ string) Stream {
	return &EmptyStream{}
}

// Info 記錄一般訊息, 用於記錄一般訊息
func (this *EmptyLogger) Info(_ string) Stream {
	return &EmptyStream{}
}

// Warn 記錄警告訊息, 用於記錄遊戲邏輯錯誤
func (this *EmptyLogger) Warn(_ string) Stream {
	return &EmptyStream{}
}

// Error 記錄錯誤訊息, 用於記錄伺服器錯誤
func (this *EmptyLogger) Error(_ string) Stream {
	return &EmptyStream{}
}

// EmptyStream 空記錄
type EmptyStream struct {
}

// Message 記錄訊息
func (this *EmptyStream) Message(_ string, _ ...any) Stream {
	return this
}

// Caller 記錄呼叫訊息
func (this *EmptyStream) Caller(_ int) Stream {
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
