package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/pools"
)

const namePool = "pool" // 特性名稱

// NewPool 建立執行緒池資料
func NewPool() *Pool {
	return &Pool{}
}

// Pool 執行緒池資料
type Pool struct {
	config pools.Config // 配置資料
}

// Initialize 初始化處理
func (this *Pool) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(namePool, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", namePool, err)
	} // if

	this.config.Logger = &poolLogger{}

	if err := mizugos.Poolmgr().Initialize(&this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", namePool, err)
	} // if

	System.Info(namePool).Caller(0).Message("initialize").KV("config", &this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Pool) Finalize() {
	mizugos.Logmgr().Finalize()
}

// poolLogger 執行緒日誌
type poolLogger struct {
}

// Printf 輸出日誌
func (this *poolLogger) Printf(format string, args ...any) {
	System.Error(namePool).Caller(1).Message(format, args...).End()
}
