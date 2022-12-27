package commons

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/logs"
)

// NewLogger 建立日誌資料
func NewLogger() *Logger {
	return &Logger{
		name: "logger",
	}
}

// Logger 日誌資料
type Logger struct {
	name string // 日誌名稱
}

// Initialize 初始化處理
func (this *Logger) Initialize() error {
	logger := logs.ZapLogger{}

	if err := mizugos.Configmgr().ReadFile(this.name, "yaml"); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &logger); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Logmgr().Initialize(&logger); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	return nil
}

// Finalize 結束處理
func (this *Logger) Finalize() {
	mizugos.Logmgr().Finalize()
}
