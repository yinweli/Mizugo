package logs

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// NewLogmgr 建立日誌管理器
func NewLogmgr() *Logmgr {
	return &Logmgr{}
}

// Logmgr 日誌管理器, 用於執行管理日誌相關功能; 使用前需要執行 Initialize, 使用完畢需要執行 Finalize
//
// 使用者可以選擇要使用以下哪種日誌
//   - EmptyLogger: 空日誌, 不會輸出任何訊息
//   - ZapLogger: uber實現的高效能日誌功能
//   - 自訂日誌: 如果使用者想要自訂日誌, 需要實現 Logger 介面與 Stream 介面
//
// 目前提供以下日誌等級: Debug, Info, Warn, Error
type Logmgr struct {
	once   utils.SyncOnce         // 單次執行物件
	logger utils.SyncAttr[Logger] // 日誌物件
}

// Initialize 初始化處理
func (this *Logmgr) Initialize(logger Logger) (err error) {
	if this.once.Done() {
		return fmt.Errorf("logmgr initialize: already initialize")
	} // if

	this.once.Do(func() {
		if logger == nil {
			logger = &EmptyLogger{} // 預設使用空日誌
		} // if

		if err = logger.Initialize(); err != nil {
			err = fmt.Errorf("logmgr initialize: %w", err)
			return
		} // if

		this.logger.Set(logger)
	})

	return err
}

// Finalize 結束處理
func (this *Logmgr) Finalize() {
	if this.once.Done() == false {
		return
	} // if

	this.logger.Get().Finalize()
}

// Debug 記錄除錯訊息
func (this *Logmgr) Debug(label string) Stream {
	return this.logger.Get().New(label, LevelDebug)
}

// Info 記錄一般訊息
func (this *Logmgr) Info(label string) Stream {
	return this.logger.Get().New(label, LevelInfo)
}

// Warn 記錄警告訊息
func (this *Logmgr) Warn(label string) Stream {
	return this.logger.Get().New(label, LevelWarn)
}

// Error 記錄錯誤訊息
func (this *Logmgr) Error(label string) Stream {
	return this.logger.Get().New(label, LevelError)
}
