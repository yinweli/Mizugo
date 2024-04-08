package loggers

import (
	"fmt"
	"sync"
)

// NewLogmgr 建立日誌管理器
func NewLogmgr() *Logmgr {
	return &Logmgr{
		logger: map[string]Logger{},
	}
}

// Logmgr 日誌管理器, 用於執行管理日誌相關功能;
// 使用前需要執行 Add 新增日誌; 使用完畢需要執行 Finalize 結束所有日誌
//
// 新增日誌時, 有以下預設日誌物件可用
//   - EmptyLogger: 空日誌, 不會輸出任何訊息
//   - ZapLogger: uber實現的高效能日誌功能
//
// 如果使用者想要自訂日誌, 需要實現 Logger 介面與 Stream 介面
//
// Logmgr 提供以下日誌函式: Debug, Info, Warn, Error
type Logmgr struct {
	logger map[string]Logger // 日誌列表
	lock   sync.RWMutex      // 執行緒鎖
}

// Add 新增日誌物件
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

// Get 取得日誌物件
func (this *Logmgr) Get(name string) Logger {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if logger, ok := this.logger[name]; ok {
		return logger
	} // if

	return nil
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
