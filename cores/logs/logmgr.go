package logs

import (
	"fmt"
	"sync"
)

// NewLogmgr 建立日誌管理器
func NewLogmgr() *Logmgr {
	return &Logmgr{
		logger: &EmptyLogger{}, // 預設使用空日誌
	}
}

// Logmgr 日誌管理器
type Logmgr struct {
	logger Logger       // 日誌物件
	lock   sync.RWMutex // 執行緒鎖
}

// Initialize 初始化處理
func (this *Logmgr) Initialize(logger Logger) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := logger.Initialize(); err != nil {
		return fmt.Errorf("logmgr initialize: %w", err)
	} // if

	this.logger = logger
	return nil
}

// Finalize 結束處理
func (this *Logmgr) Finalize() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.logger.Finalize()
}

// Debug 記錄除錯訊息
func (this *Logmgr) Debug(label string) Stream {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.logger.New(label, LevelDebug)
}

// Info 記錄一般訊息
func (this *Logmgr) Info(label string) Stream {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.logger.New(label, LevelInfo)
}

// Warn 記錄警告訊息
func (this *Logmgr) Warn(label string) Stream {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.logger.New(label, LevelWarn)
}

// Error 記錄錯誤訊息
func (this *Logmgr) Error(label string) Stream {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.logger.New(label, LevelError)
}
