package loggers

import (
	"fmt"
)

// EmptyLogger 空日誌, 不會輸出任何訊息
type EmptyLogger struct {
	fail bool // 提供給單元測試使用的初始化旗標
}

// Initialize 初始化處理
func (this *EmptyLogger) Initialize() error {
	if this.fail == false {
		return nil
	} else {
		return fmt.Errorf("initialize failed")
	} // if
}

// Finalize 結束處理
func (this *EmptyLogger) Finalize() {
}

// Get 取得 Retain 儲存器
func (this *EmptyLogger) Get() Retain {
	return emptyRetain
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

// Debug 建立除錯訊息的 Stream
func (this *EmptyRetain) Debug(_ string) Stream {
	return emptyStream
}

// Info 建立一般訊息的 Stream
func (this *EmptyRetain) Info(_ string) Stream {
	return emptyStream
}

// Warn 建立警告訊息的 Stream
func (this *EmptyRetain) Warn(_ string) Stream {
	return emptyStream
}

// Error 建立錯誤訊息的 Stream
func (this *EmptyRetain) Error(_ string) Stream {
	return emptyStream
}

// EmptyStream 空記錄
type EmptyStream struct {
}

// Message 記錄文字訊息
func (this *EmptyStream) Message(_ string, _ ...any) Stream {
	return this
}

// KV 記錄鍵值訊息
func (this *EmptyStream) KV(_ string, _ any) Stream {
	return this
}

// Caller 記錄呼叫位置
func (this *EmptyStream) Caller(_ int, _ ...bool) Stream {
	return this
}

// Error 記錄錯誤物件
func (this *EmptyStream) Error(_ error) Stream {
	return this
}

// End 結束記錄, 將記錄交回 Retain
func (this *EmptyStream) End() Retain {
	return emptyRetain
}

// EndFlush 結束記錄, 將記錄交回 Retain 並立即儲存
func (this *EmptyStream) EndFlush() {
}

var emptyRetain = &EmptyRetain{}
var emptyStream = &EmptyStream{}
