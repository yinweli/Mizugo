package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/entrys"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
)

func main() {
	mizugos.Start("example_clientgo", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	client.logger = features.NewLogger()
	client.pool = features.NewPool()
	client.metrics = features.NewMetrics()
	client.pingJson = entrys.NewPingJson()
	client.pingProto = entrys.NewPingProto()
	client.pingStack = entrys.NewPingStack()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := client.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.pool.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.pingJson.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.pingProto.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.pingStack.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	client.pingJson.Finalize()
	client.pingProto.Finalize()
	client.pingStack.Finalize()
	client.metrics.Finalize()
	client.pool.Finalize()
	client.logger.Finalize()
}

// client 客戶端資料
var client struct {
	logger    *features.Logger  // 日誌資料
	pool      *features.Pool    // 執行緒池資料
	metrics   *features.Metrics // 統計資料
	pingJson  *entrys.PingJson  // PingJson入口
	pingProto *entrys.PingProto // PingProto入口
	pingStack *entrys.PingStack // PingStack入口
}
