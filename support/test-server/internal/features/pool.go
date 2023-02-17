package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/pools"
)

// NewPool 建立執行緒池資料
func NewPool() *Pool {
	return &Pool{
		name: "pool",
	}
}

// Pool 執行緒池資料
type Pool struct {
	name   string       // 特性名稱
	config pools.Config // 配置資料
}

// Initialize 初始化處理
func (this *Pool) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	this.config.Logger = &poolLogger{}

	if err := mizugos.Poolmgr().Initialize(&this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Info(this.name).Message("initialize").KV("config", &this.config).End()
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
	mizugos.Logmgr().Error("pool").Message(format, args...).End()
}
