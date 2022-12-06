package logs

// 提供日誌相關功能, 可以自選執行的日誌機制
// 初始化日誌: 執行logs.Initialize(日誌物件)
// 結束日誌: 執行logs.Finalize()
// 記錄日誌: 執行以下之一來取得記錄物件, 再通過記錄介面來寫入資訊
//     logs.Debug(日誌標籤)
//     logs.Info(日誌標籤)
//     logs.Warn(日誌標籤)
//     logs.Error(日誌標籤)

const (
	LevelDebug Level = iota // 除錯訊息
	LevelInfo               // 一般訊息
	LevelWarn               // 警告訊息
	LevelError              // 錯誤訊息
)

// Logger 日誌介面, 實作時需要注意會在多執行緒環境下運作
type Logger interface {
	// Initialize 初始化處理
	Initialize()

	// Finalize 結束處理
	Finalize()

	// New 建立日誌
	New(label string, level Level) Stream
}

// Stream 記錄介面, 實作時需要注意會在多執行緒環境下運作
type Stream interface {
	// Message 記錄訊息
	Message(message string) Stream

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

// Initialize 初始化日誌
func Initialize(logger Logger) {
	if logger != nil {
		logger_ = logger
	} else {
		logger_ = &EmptyLogger{}
	} // if

	logger_.Initialize()
}

// Finalize 結束日誌
func Finalize() {
	logger_.Finalize()
}

// Debug 記錄除錯訊息
func Debug(label string) Stream {
	return logger_.New(label, LevelDebug)
}

// Info 記錄一般訊息
func Info(label string) Stream {
	return logger_.New(label, LevelInfo)
}

// Warn 記錄警告訊息
func Warn(label string) Stream {
	return logger_.New(label, LevelWarn)
}

// Error 記錄錯誤訊息
func Error(label string) Stream {
	return logger_.New(label, LevelError)
}

var logger_ Logger = &EmptyLogger{} // 日誌物件
