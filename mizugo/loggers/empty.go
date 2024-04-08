package loggers

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

// Get 取得儲存器
func (this *EmptyLogger) Get() Retain {
	return &EmptyRetain{}
}

// EmptyRetain 空儲存
type EmptyRetain struct {
}

// Clear 清空內部 Stream 列表
func (this *EmptyRetain) Clear() Retain {
	return this
}

// Flush 儲存並清空內部 Stream 列表
func (this *EmptyRetain) Flush() Retain {
	return this
}

// Debug 記錄除錯訊息, 用於記錄除錯訊息
func (this *EmptyRetain) Debug(_ string) Stream {
	return &EmptyStream{retain: this}
}

// Info 記錄一般訊息, 用於記錄一般訊息
func (this *EmptyRetain) Info(_ string) Stream {
	return &EmptyStream{retain: this}
}

// Warn 記錄警告訊息, 用於記錄邏輯錯誤
func (this *EmptyRetain) Warn(_ string) Stream {
	return &EmptyStream{retain: this}
}

// Error 記錄錯誤訊息, 用於記錄嚴重錯誤
func (this *EmptyRetain) Error(_ string) Stream {
	return &EmptyStream{retain: this}
}

// EmptyStream 空記錄
type EmptyStream struct {
	retain Retain
}

// Message 記錄訊息
func (this *EmptyStream) Message(_ string, _ ...any) Stream {
	return this
}

// KV 記錄索引與數值
func (this *EmptyStream) KV(_ string, _ any) Stream {
	return this
}

// Caller 記錄呼叫位置
func (this *EmptyStream) Caller(_ int) Stream {
	return this
}

// Error 記錄錯誤
func (this *EmptyStream) Error(_ error) Stream {
	return this
}

// End 結束記錄, 並把記錄加回到 Retain 中
func (this *EmptyStream) End() Retain {
	return this.retain
}

// EndFlush 結束記錄, 並把記錄加回到 Retain 中, 然後儲存記錄
func (this *EmptyStream) EndFlush() {
}
