package logs

const (
	LevelDebug Level = iota // 除錯訊息
	LevelInfo               // 一般訊息
	LevelWarn               // 警告訊息
	LevelError              // 錯誤訊息
)

// Logger 日誌介面, 實作時需要注意會在多執行緒環境下運作
type Logger interface {
	// Initialize 初始化處理
	Initialize() error

	// Finalize 結束處理
	Finalize()

	// New 建立日誌
	New(label string, level Level) Stream
}

// Stream 記錄介面, 實作時需要注意會在多執行緒環境下運作
type Stream interface {
	// Message 記錄訊息
	Message(format string, a ...any) Stream

	// KV 記錄索引與數值
	KV(key string, value any) Stream

	// Error 記錄錯誤
	Error(err error) Stream

	// EndError 以錯誤結束記錄
	EndError(err error) error

	// End 結束記錄
	End()
}

// Level 日誌等級
type Level int
