package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/entrys"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
)

func main() {
	mizugos.Start("example_client_go", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	client.logger = features.NewLogger()
	client.pool = features.NewPool()
	client.metrics = features.NewMetrics()
	client.echo = entrys.NewEcho()
	client.ping = entrys.NewPing()

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

	if err := client.echo.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.ping.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	client.echo.Finalize()
	client.ping.Finalize()
	client.metrics.Finalize()
	client.pool.Finalize()
	client.logger.Finalize()
}

// client 客戶端資料
var client struct {
	logger  *features.Logger  // 日誌資料
	pool    *features.Pool    // 執行緒池資料
	metrics *features.Metrics // 統計資料
	echo    *entrys.Echo      // 回音入口
	ping    *entrys.Ping      // Ping入口
}
