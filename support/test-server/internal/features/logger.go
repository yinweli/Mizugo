package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos"
	"github.com/yinweli/Mizugo/v2/mizugos/loggers"
)

// LoggerInitialize 初始化日誌
func LoggerInitialize() (err error) {
	if LogCrash, err = newLogger("log-crash"); err != nil {
		return fmt.Errorf("logger initialize: %w", err)
	} // if

	if LogSystem, err = newLogger("log-system"); err != nil {
		return fmt.Errorf("logger initialize: %w", err)
	} // if

	LogSystem.Get().Info("logger").Message("initialize").EndFlush()
	return nil
}

// newLogger 建立日誌物件
func newLogger(name string) (logger loggers.Logger, err error) {
	logger = &loggers.ZapLogger{}

	if err = mizugos.Config.Unmarshal(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	if err = mizugos.Logger.Add(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	return logger, nil
}

var LogCrash loggers.Logger  // 崩潰日誌物件
var LogSystem loggers.Logger // 系統日誌物件
