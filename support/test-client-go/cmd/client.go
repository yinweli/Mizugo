package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/entrys"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
)

func main() {
	mizugos.Start(defines.CmdClient, initialize, finalize, defines.Crashlize)
}

// initialize 初始化處理
func initialize() error {
	server.logger = features.NewLogger()
	server.pool = features.NewPool()
	server.metrics = features.NewMetrics()
	server.auth = entrys.NewAuth()
	server.json = entrys.NewJson()
	server.proto = entrys.NewProto()

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

	if err := server.auth.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.json.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.proto.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	server.auth.Finalize()
	server.json.Finalize()
	server.proto.Finalize()
	server.metrics.Finalize()
	server.pool.Finalize()
	server.logger.Finalize()
}

// server 伺服器資料
var server struct {
	logger  *features.Logger  // 日誌資料
	pool    *features.Pool    // 執行緒池資料
	metrics *features.Metrics // 統計資料
	auth    *entrys.Auth      // Auth入口
	json    *entrys.Json      // Json入口
	proto   *entrys.Proto     // Proto入口
}