package commons

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/logs"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
)

// NewLogger 建立日誌資料
func NewLogger() *Logger {
	return &Logger{
		name: "logger",
	}
}

// Logger 日誌資料
type Logger struct {
	name   string         // 日誌名稱
	config logs.ZapLogger // 設定資料
}

// Initialize 初始化處理
func (this *Logger) Initialize() error {
	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Logmgr().Initialize(&this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Info(this.name).Message("initialize").KV("config", &this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Logger) Finalize() {
	mizugos.Logmgr().Finalize()
}
