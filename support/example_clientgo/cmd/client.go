package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/commons"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/entrys"
)

func main() {
	mizugos.Start("example_client_go", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	client.logger = commons.NewLogger()
	client.metrics = commons.NewMetrics()
	client.echo = entrys.NewEcho()
	client.ping = entrys.NewPing()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := client.logger.Initialize(); err != nil {
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
	client.logger.Finalize() // 日誌必須最後結束
}

// client 客戶端資料
var client struct {
	logger  *commons.Logger  // 日誌資料
	metrics *commons.Metrics // 統計資料
	echo    *entrys.Echo     // 回音入口
	ping    *entrys.Ping     // Ping入口
}
