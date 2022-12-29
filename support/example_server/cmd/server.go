package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/features/commons"
	"github.com/yinweli/Mizugo/support/example_server/features/defines"
	"github.com/yinweli/Mizugo/support/example_server/features/entrys"
)

func main() {
	mizugos.Start("example_server", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	feature.logger = commons.NewLogger()
	feature.echo = entrys.NewEcho()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := feature.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.echo.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	feature.echo.Finalize()
	feature.metrics.Finalize()
	feature.logger.Finalize() // 日誌必須最後結束
}

// feature 功能資料
var feature struct {
	logger  *commons.Logger // 日誌資料
	metrics *entrys.Metrics // 統計入口
	echo    *entrys.Echo    // 回音入口
}
