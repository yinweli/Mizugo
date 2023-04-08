package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/logs"
)

// NewLogger 建立日誌資料
func NewLogger() *Logger {
	return &Logger{
		name: "feature logger",
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

	LogSystem.Info(this.name).Caller(0).Message("initialize crash log").KV("config", LogCrash).End()
	LogSystem.Info(this.name).Caller(0).Message("initialize system log").KV("config", LogSystem).End()
	return nil
}

// Finalize 結束處理
func (this *Logger) Finalize() {
	mizugos.Logmgr().Finalize()
}

// create 建立日誌物件
func (this *Logger) create(name string) (logger logs.Logger, err error) {
	logger = &logs.ZapLogger{}

	if err = mizugos.Configmgr().Unmarshal(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	if err = mizugos.Logmgr().Add(name, logger); err != nil {
		return nil, fmt.Errorf("create: %v: %w", name, err)
	} // if

	return logger, nil
}

var LogCrash logs.Logger  // 崩潰日誌物件
var LogSystem logs.Logger // 系統日誌物件
