package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/internal/commons"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/entrys"
)

func main() {
	mizugos.Start("example_server", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	server.logger = commons.NewLogger()
	server.metrics = commons.NewMetrics()
	server.echo = entrys.NewEcho()
	server.ping = entrys.NewPing()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := server.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.echo.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.ping.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	server.echo.Finalize()
	server.ping.Finalize()
	server.metrics.Finalize()
	server.logger.Finalize() // 日誌必須最後結束
}

// server 伺服器資料
var server struct {
	logger  *commons.Logger  // 日誌資料
	metrics *commons.Metrics // 統計資料
	echo    *entrys.Echo     // 回音入口
	ping    *entrys.Ping     // Ping入口
}
