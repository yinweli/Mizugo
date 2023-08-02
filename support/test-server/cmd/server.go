package main

import (
	"fmt"
	"runtime/debug"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/entrys"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
)

func main() {
	defer func() {
		if cause := recover(); cause != nil {
			features.LogCrash.Get().Error("crash").KV("stack", string(debug.Stack())).Error(fmt.Errorf("%s", cause)).EndFlush()
		} // if
	}()

	ctx := ctxs.Get().WithCancel()
	name := defines.CmdServer
	mizugos.Start()

	if err := initialize(); err != nil {
		fmt.Println(fmt.Errorf("%v start: %w", name, err))
		mizugos.Stop()
		return
	} // if

	fmt.Printf("%v start\n", name)

	for range ctx.Done() {
		// do nothing...
	} // for

	finalize()
	mizugos.Stop()
	fmt.Printf("%v shutdown\n", name)
}

// initialize 初始化處理
func initialize() error {
	server.logger = features.NewLogger()
	server.pool = features.NewPool()
	server.metrics = features.NewMetrics()
	server.redmo = features.NewRedmo()
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

	if err := server.redmo.Initialize(); err != nil {
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
	server.redmo.Finalize()
	server.metrics.Finalize()
	server.pool.Finalize()
	server.logger.Finalize()
}

// server 伺服器資料
var server struct {
	logger  *features.Logger  // 日誌資料
	pool    *features.Pool    // 執行緒池資料
	metrics *features.Metrics // 統計資料
	redmo   *features.Redmo   // 資料庫資料
	auth    *entrys.Auth      // Auth入口
	json    *entrys.Json      // Json入口
	proto   *entrys.Proto     // Proto入口
}
