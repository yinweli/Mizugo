package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugo/loggers"
)

// InitializeLogger 初始化日誌管理器
func InitializeLogger() (err error) {
	name := "logger"
	Logger = loggers.NewLogmgr()

	if LogCrash, err = newLogger("log-crash"); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	if LogSystem, err = newLogger("log-system"); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

// FinalizeLogger 結束日誌管理器
func FinalizeLogger() {
	if Logger != nil {
		Logger.Finalize()
	} // if
}

// newLogger 建立日誌物件
func newLogger(name string) (logger loggers.Logger, err error) {
	logger = &loggers.ZapLogger{}

	if err = Config.Unmarshal(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	if err = Logger.Add(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	return logger, nil
}

var Logger *loggers.Logmgr   // 日誌管理器
var LogCrash loggers.Logger  // 崩潰日誌物件
var LogSystem loggers.Logger // 系統日誌物件
