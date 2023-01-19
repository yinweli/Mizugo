package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/entrys"
	"github.com/yinweli/Mizugo/support/example_server/internal/features"
)

func main() {
	mizugos.Start("example_server", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	server.logger = features.NewLogger()
	server.pool = features.NewPool()
	server.metrics = features.NewMetrics()
	server.pingJson = entrys.NewPingJson()
	server.pingProto = entrys.NewPingProto()
	server.pingStack = entrys.NewPingStack()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := server.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.pool.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.pingJson.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.pingProto.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.pingStack.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	server.pingJson.Finalize()
	server.pingProto.Finalize()
	server.pingStack.Finalize()
	server.metrics.Finalize()
	server.pool.Finalize()
	server.logger.Finalize()
}

// server 伺服器資料
var server struct {
	logger    *features.Logger  // 日誌資料
	pool      *features.Pool    // 執行緒池資料
	metrics   *features.Metrics // 統計資料
	pingJson  *entrys.PingJson  // PingJson入口
	pingProto *entrys.PingProto // PingProto入口
	pingStack *entrys.PingStack // PingStack入口
}
