package logs

const (
	LevelDebug = "debug" // 除錯訊息
	LevelInfo  = "info"  // 一般訊息
	LevelWarn  = "warn"  // 警告訊息
	LevelError = "error" // 錯誤訊息
)

// Logger 日誌介面, 實作時需要注意會在多執行緒環境下運作
type Logger interface {
	// Initialize 初始化處理
	Initialize() error

	// Finalize 結束處理
	Finalize()

	// Debug 記錄除錯訊息, 用於記錄除錯訊息
	Debug(label string) Stream

	// Info 記錄一般訊息, 用於記錄一般訊息
	Info(label string) Stream

	// Warn 記錄警告訊息, 用於記錄遊戲邏輯錯誤
	Warn(label string) Stream

	// Error 記錄錯誤訊息, 用於記錄伺服器錯誤
	Error(label string) Stream
}

// Stream 記錄介面, 實作時需要注意會在多執行緒環境下運作
type Stream interface {
	// Message 記錄訊息
	Message(format string, a ...any) Stream

	// Caller 記錄呼叫訊息
	Caller(skip int) Stream

	// KV 記錄索引與數值
	KV(key string, value any) Stream

	// Error 記錄錯誤
	Error(err error) Stream

	// EndError 以錯誤結束記錄
	EndError(err error)

	// End 結束記錄
	End()
}
