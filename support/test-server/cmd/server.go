package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/entrys"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
)

func main() {
	mizugos.Start("test_server", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	server.logger = features.NewLogger()
	server.pool = features.NewPool()
	server.metrics = features.NewMetrics()
	server.Json = entrys.NewJson()
	server.Proto = entrys.NewProto()
	server.Stack = entrys.NewStack()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := mizugos.Configmgr().ReadFile(defines.ConfigFile, defines.ConfigType); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.pool.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.Json.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.Proto.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.Stack.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	server.Json.Finalize()
	server.Proto.Finalize()
	server.Stack.Finalize()
	server.metrics.Finalize()
	server.pool.Finalize()
	server.logger.Finalize()
}

// server 伺服器資料
var server struct {
	logger  *features.Logger  // 日誌資料
	pool    *features.Pool    // 執行緒池資料
	metrics *features.Metrics // 統計資料
	Json    *entrys.Json      // Json入口
	Proto   *entrys.Proto     // Proto入口
	Stack   *entrys.Stack     // Stack入口
}
