package echos

import (
	"fmt"
	"path/filepath"

	"github.com/yinweli/Mizugo/mizugos"
)

// 這裡會執行使用字串以及md5編碼作為封包內容的回音伺服器

const serverName = "serv_echo" // 伺服器名稱
const configFile = "echo.yaml" // 伺服器設定檔案

// config 伺服器設定資料
type config struct {
}

// Initialize 初始化處理
func Initialize(configPath string) error {
	if err := mizugos.Configmgr().ReadFile(filepath.Join(configPath, configFile)); err != nil {
		return fmt.Errorf("%v initialize: %w", serverName, err)
	} // if

	return nil
}

// Finalize 結束處理
func Finalize() {

}
