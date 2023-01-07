package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/commons"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
)

func main() {
	mizugos.Start("example_client_go", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	feature.logger = commons.NewLogger()
	feature.metrics = commons.NewMetrics()
	feature.echoSingle = features.NewEchoSingle()
	feature.echoCycle = features.NewEchoCycle()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := feature.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.metrics.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.echoSingle.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.echoCycle.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	feature.echoSingle.Finalize()
	feature.echoCycle.Finalize()
	feature.metrics.Finalize()
	feature.logger.Finalize() // 日誌必須最後結束
}

// feature 功能資料
var feature struct {
	logger     *commons.Logger      // 日誌資料
	metrics    *commons.Metrics     // 統計資料
	echoSingle *features.EchoSingle // 單次回音
	echoCycle  *features.EchoCycle  // 循環回音
}
