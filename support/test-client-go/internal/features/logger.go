package features

import (
	"fmt"
	"runtime/debug"

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

	if System = mizugos.Logmgr().Get(defines.LogSystem); System == nil {
		return fmt.Errorf("%v initialize: system logger nil", nameLogger)
	} // if

	if err := this.add(defines.LogCrash, &this.configCrash); err != nil {
		return fmt.Errorf("%v initialize: %w", nameLogger, err)
	} // if

	if Crash = mizugos.Logmgr().Get(defines.LogCrash); Crash == nil {
		return fmt.Errorf("%v initialize: crash logger nil", nameLogger)
	} // if

	System.Info(nameLogger).Caller(0).Message("initialize system log").KV("config", &this.configSystem).End()
	System.Info(nameLogger).Caller(0).Message("initialize crash log").KV("config", &this.configCrash).End()
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

// Crashlize 崩潰處理
func Crashlize(cause any) {
	Crash.Error("crash").KV("stack", string(debug.Stack())).EndError(fmt.Errorf("%s", cause))
}

var System logs.Logger // 系統日誌
var Crash logs.Logger  // 崩潰日誌
