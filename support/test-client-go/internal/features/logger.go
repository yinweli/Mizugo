package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/loggers"
)

// NewLogger 建立日誌資料
func NewLogger() *Logger {
	return &Logger{
		name: "logger",
	}
}

// Logger 日誌資料
type Logger struct {
	name string // 系統名稱
}

// Initialize 初始化處理
func (this *Logger) Initialize() (err error) {
	if LogCrash, err = this.create("log-crash"); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if LogSystem, err = this.create("log-system"); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	LogSystem.Get().
		Info(this.name).Message("initialize crash log").KV("config", LogCrash).Caller(0).End().
		Info(this.name).Message("initialize system log").KV("config", LogSystem).Caller(0).End().
		Flush()
	return nil
}

// Finalize 結束處理
func (this *Logger) Finalize() {
	mizugos.Logmgr().Finalize()
}

// create 建立日誌物件
func (this *Logger) create(name string) (logger loggers.Logger, err error) {
	logger = &loggers.ZapLogger{}

	if err = mizugos.Configmgr().Unmarshal(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	if err = mizugos.Logmgr().Add(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	return logger, nil
}

var LogCrash loggers.Logger  // 崩潰日誌物件
var LogSystem loggers.Logger // 系統日誌物件
