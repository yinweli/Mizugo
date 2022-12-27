package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/commons"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/entrys"
)

func main() {
	mizugos.Start("example_client_go", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	feature.logger = commons.NewLogger()
	feature.echoc = entrys.NewEchoc()

	mizugos.Configmgr().AddPath(defines.ConfigPath)

	if err := feature.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.echoc.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	feature.logger.Finalize()
	feature.echoc.Finalize()
}

// feature 功能資料
var feature struct {
	logger *commons.Logger // 日誌資料
	echoc  *entrys.Echoc   // 回音入口
}