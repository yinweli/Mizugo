package loggers

import (
	"fmt"
	"sync"
)

const (
	LevelDebug = "debug" // 除錯訊息
	LevelInfo  = "info"  // 一般訊息
	LevelWarn  = "warn"  // 警告訊息
	LevelError = "error" // 錯誤訊息
)

// NewLogmgr 建立日誌管理器
func NewLogmgr() *Logmgr {
	return &Logmgr{
		logger: map[string]Logger{},
	}
}

// Logmgr 日誌管理器, 負責日誌實例的註冊, 取得與銷毀
//
// 使用方式:
//   - 呼叫 Add() 新增日誌實例
//   - 透過 Get() 取得日誌並記錄訊息
//   - 程式結束時呼叫 Finalize() 釋放資源
//
// 內建日誌實作:
//   - EmptyLogger: 空實作, 不輸出任何訊息
//   - ZapLogger: 基於 uber/zap 的高效能日誌
//
// 若要自訂日誌, 請實作 Logger 與其對應的 Retain / Stream
type Logmgr struct {
	logger map[string]Logger // 日誌列表
	lock   sync.RWMutex      // 執行緒鎖
}

// Add 新增日誌實例並初始化, name 不可為空且不可重複
func (this *Logmgr) Add(name string, logger Logger) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if name == "" {
		return fmt.Errorf("logmgr add: name empty")
	} // if

	if _, ok := this.logger[name]; ok {
		return fmt.Errorf("logmgr add: name duplicate")
	} // if

	if logger == nil {
		return fmt.Errorf("logmgr add: logger nil")
	} // if

	if err := logger.Initialize(); err != nil {
		return fmt.Errorf("logmgr add: %w", err)
	} // if

	this.logger[name] = logger
	return nil
}

// Get 以名稱取得日誌
func (this *Logmgr) Get(name string) Logger {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.logger[name]
}

// Finalize 結束處理
func (this *Logmgr) Finalize() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.logger {
		itor.Finalize()
	} // for

	this.logger = map[string]Logger{}
}

// Logger 日誌介面
//   - 必須具備多執行緒安全性
//   - 允許被多個執行緒同時共用
type Logger interface {
	// Initialize 初始化處理
	Initialize() error

	// Finalize 結束處理
	Finalize()

	// Get 取得 Retain 儲存器
	Get() Retain
}

// Retain 儲存介面
//   - 用於暫存多筆 Stream, 再一次性送出
//   - 方法可能在多執行緒環境下被呼叫, 但同一個 Retain 實例不可跨執行緒共用
//   - 每個執行緒應建立並使用自己的 Retain
type Retain interface {
	// Clear 清空內部 Stream 列表
	Clear() Retain

	// Flush 儲存並清空內部 Stream 列表
	Flush() Retain

	// Debug 建立除錯訊息的 Stream
	Debug(label string) Stream

	// Info 建立一般訊息的 Stream
	Info(label string) Stream

	// Warn 建立警告訊息的 Stream
	Warn(label string) Stream

	// Error 建立錯誤訊息的 Stream
	Error(label string) Stream
}

// Stream 記錄介面
//   - 用於記錄單筆訊息的細節
//   - 方法可能在多執行緒環境下被呼叫, 但同一個 Stream 實例不可跨執行緒共用
//   - 每個執行緒應建立並使用自己的 Stream
type Stream interface {
	// Message 記錄文字訊息
	Message(format string, a ...any) Stream

	// KV 記錄鍵值訊息
	KV(key string, value any) Stream

	// Caller 記錄呼叫位置
	//   - skip: 呼叫堆疊跳過層數
	//   - simple: 是否輸出精簡函式名稱
	Caller(skip int, simple ...bool) Stream

	// Error 記錄錯誤物件
	Error(err error) Stream

	// End 結束記錄, 將記錄交回 Retain
	End() Retain

	// EndFlush 結束記錄, 將記錄交回 Retain 並立即儲存
	EndFlush()
}
