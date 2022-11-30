package logs

// Level 日誌等級
type Level int

const (
	LevelDebug Level = iota // 除錯訊息
	LevelInfo               // 一般訊息
	LevelWarn               // 警告訊息
	LevelError              // 錯誤訊息
)

// Logger 日誌介面, 實作時需要注意會在多執行緒環境下運作
type Logger interface {
	// Message 記錄訊息
	Message(message string) Logger

	// KV 記錄索引與數值
	KV(key string, value any) Logger

	// Error 記錄錯誤
	Error(err error) Logger

	// EndError 以錯誤結束記錄
	EndError(err error) error

	// End 結束記錄
	End()
}

// NewLog 建立日誌函式類型
type NewLog func(label string, level Level) Logger

// Set 設定建立日誌函式
func Set(f NewLog) {
	if f != nil {
		newLog = f
	} else {
		newLog = NewEmpty
	} // if
}

// Debug 記錄除錯訊息
func Debug(label string) Logger {
	return newLog(label, LevelDebug)
}

// Info 記錄一般訊息
func Info(label string) Logger {
	return newLog(label, LevelInfo)
}

// Warn 記錄警告訊息
func Warn(label string) Logger {
	return newLog(label, LevelWarn)
}

// Error 記錄錯誤訊息
func Error(label string) Logger {
	return newLog(label, LevelError)
}

var newLog NewLog = NewEmpty // 建立日誌函式
