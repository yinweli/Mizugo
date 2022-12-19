package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/server/commons"
	"github.com/yinweli/Mizugo/support/example_server/server/echos"
)

const configPath = "config" // 設定檔案路徑

func main() {
	mizugos.Start("example", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	serv.logger = commons.NewLogger(configPath)
	serv.echoServer = echos.NewServer(configPath)

	if err := serv.logger.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err := serv.echoServer.Initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	serv.echoServer.Finalize()
	serv.logger.Finalize()
}

var serv struct {
	logger     *commons.Logger // 日誌資料
	echoServer *echos.Server   // 回音伺服器資料
}
