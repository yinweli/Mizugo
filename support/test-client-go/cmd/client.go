package main

import (
	"fmt"
	"runtime/debug"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/entrys"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
)

func main() {
	defer func() {
		if cause := recover(); cause != nil {
			features.LogCrash.Get().Error("crash").KV("stack", string(debug.Stack())).Error(fmt.Errorf("%s", cause)).EndFlush()
		} // if
	}()

	ctx := ctxs.Get().WithCancel()
	name := defines.CmdClient
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
	client.logger = features.NewLogger()
	client.pool = features.NewPool()
	client.metrics = features.NewMetrics()
	client.auth = entrys.NewAuth()
	client.json = entrys.NewJson()
	client.proto = entrys.NewProto()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := mizugos.Configmgr().ReadFile(defines.ConfigFile, defines.ConfigType); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.pool.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.auth.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.json.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := client.proto.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	client.auth.Finalize()
	client.json.Finalize()
	client.proto.Finalize()
	client.metrics.Finalize()
	client.pool.Finalize()
	client.logger.Finalize()
}

// client 客戶端資料
var client struct {
	logger  *features.Logger  // 日誌資料
	pool    *features.Pool    // 執行緒池資料
	metrics *features.Metrics // 度量資料
	auth    *entrys.Auth      // Auth入口
	json    *entrys.Json      // Json入口
	proto   *entrys.Proto     // Proto入口
}
