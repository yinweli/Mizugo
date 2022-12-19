package commons

import (
	"fmt"
	"path/filepath"

	"github.com/yinweli/Mizugo/cores/logs"
	"github.com/yinweli/Mizugo/mizugos"
)

// NewLogger 建立日誌資料
func NewLogger(configPath string) *Logger {
	return &Logger{
		loggerName: "zapLogger",
		configName: "logger",
		configFile: "logger.yaml",
		configPath: configPath,
	}
}

// Logger 日誌資料
type Logger struct {
	loggerName string         // 日誌名稱
	configName string         // 設定名稱
	configFile string         // 設定檔案名稱
	configPath string         // 設定檔案路徑
	logger     logs.ZapLogger // 日誌物件
}

// Initialize 初始化處理
func (this *Logger) Initialize() error {
	if err := mizugos.Configmgr().ReadFile(filepath.Join(this.configPath, this.configFile)); err != nil {
		return fmt.Errorf("%v initialize: %w", this.loggerName, err)
	} // if

	if err := mizugos.Configmgr().GetObject(this.configName, &this.logger); err != nil {
		return fmt.Errorf("%v initialize: %w", this.loggerName, err)
	} // if

	if err := mizugos.Logmgr().Initialize(&this.logger); err != nil {
		return fmt.Errorf("%v initialize: %w", this.loggerName, err)
	} // if

	return nil
}

// Finalize 結束處理
func (this *Logger) Finalize() {
	mizugos.Logmgr().Finalize()
}
