package logs

var New func() Logger = nil // 建立日誌函式

// Begin 開始日誌
func Begin() Logger {
	if New != nil {
		return New()
	} // if

	return &empty{}
}

// Logger 日誌介面, 實作時需要注意會在多執行緒環境下運作
type Logger interface {
	// Message 記錄訊息
	Message(message string) Logger

	// KeyValue 記錄索引與數值
	KeyValue(key string, value any) Logger

	// Error 記錄錯誤
	Error(err error) Logger

	// End 結束記錄
	End()

	// EndError 以錯誤結束記錄
	EndError(err error) error
}

// empty 空日誌
type empty struct {
}

// Message 記錄訊息
func (this *empty) Message(message string) Logger {
	return this
}

// KeyValue 記錄索引與數值
func (this *empty) KeyValue(key string, value any) Logger {
	return this
}

// Error 記錄錯誤
func (this *empty) Error(err error) Logger {
	return this
}

// End 結束記錄
func (this *empty) End() {
}

// EndError 以錯誤結束記錄
func (this *empty) EndError(err error) error {
	return err
}
