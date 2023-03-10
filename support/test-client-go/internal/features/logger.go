package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/logs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
)

const nameLogger = "logger" // 特性名稱

// NewLogger 建立日誌資料
func NewLogger() *Logger {
	return &Logger{}
}

// Logger 日誌資料
type Logger struct {
	configSystem logs.ZapLogger // 系統日誌配置資料
	configCrash  logs.ZapLogger // 崩潰日誌配置資料
}

// Initialize 初始化處理
func (this *Logger) Initialize() error {
	if err := this.add(defines.LogSystem, &this.configSystem); err != nil {
		return fmt.Errorf("%v initialize: %w", nameLogger, err)
	} // if

	if err := this.add(defines.LogCrash, &this.configCrash); err != nil {
		return fmt.Errorf("%v initialize: %w", nameLogger, err)
	} // if

	mizugos.Info(defines.LogSystem, nameLogger).Caller(0).Message("initialize").KV("config system", &this.configSystem).End()
	mizugos.Info(defines.LogSystem, nameLogger).Caller(0).Message("initialize").KV("config crash", &this.configCrash).End()
	return nil
}

// Finalize 結束處理
func (this *Logger) Finalize() {
	mizugos.Logmgr().Finalize()
}

// add 新增日誌物件
func (this *Logger) add(name string, logger logs.Logger) error {
	if err := mizugos.Configmgr().Unmarshal(name, logger); err != nil {
		return fmt.Errorf("add: %v: %w", name, err)
	} // if

	if err := mizugos.Logmgr().Add(name, logger); err != nil {
		return fmt.Errorf("add: %v: %w", name, err)
	} // if

	return nil
}
