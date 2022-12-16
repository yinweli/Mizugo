package main

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_server/server/echos"
)

const configPath = "config" // 設定檔案路徑

func main() {
	mizugos.Start("example", initialize, finalize)
}

// initialize 初始化處理
func initialize() error {
	if err := echos.Initialize(configPath); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	return nil
}

// finalize 結束處理
func finalize() {
	echos.Finalize()
}
