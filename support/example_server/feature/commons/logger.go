package commons

import (
	"fmt"
	"path/filepath"

	"github.com/yinweli/Mizugo/cores/logs"
	"github.com/yinweli/Mizugo/mizugos"
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
func (this *Logger) Initialize(configPath string) error {
	if err := mizugos.Configmgr().ReadFile(filepath.Join(configPath, this.name+".yaml")); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	logger := logs.ZapLogger{}

	if err := mizugos.Configmgr().GetObject(this.name, &logger); err != nil {
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
