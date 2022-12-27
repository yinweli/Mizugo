package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/features/commons"
	"github.com/yinweli/Mizugo/support/example_server/features/entryechos"
)

func main() {
	mizugos.Start("example_server", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	feature.logger = commons.NewLogger()
	feature.entryEcho = entryechos.NewEntry()

	mizugos.Configmgr().AddPath("config")

	if err := feature.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := feature.entryEcho.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	feature.entryEcho.Finalize()
	feature.logger.Finalize()
}

// feature 功能資料
var feature struct {
	logger    *commons.Logger   // 日誌資料
	entryEcho *entryechos.Entry // 回音入口
}
