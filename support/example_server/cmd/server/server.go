package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/server/commons"
	"github.com/yinweli/Mizugo/support/example_server/server/entrys"
)

const configPath = "config" // 設定檔案路徑

func main() {
	mizugos.Start("example", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	server.logger = commons.NewLogger()
	server.entry.echo = entrys.NewEcho()

	if err := server.logger.Initialize(configPath); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := server.entry.echo.Initialize(configPath); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	server.entry.echo.Finalize()
	server.logger.Finalize()
}

// server 伺服器資料
var server struct {
	logger *commons.Logger // 日誌資料
	entry  struct {        // 入口資料
		echo *entrys.Echo // 回音入口資料
	}
}
