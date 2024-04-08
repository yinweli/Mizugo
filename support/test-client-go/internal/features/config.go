package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugo/configs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
)

// InitializeConfig 初始化配置管理器
func InitializeConfig() (err error) {
	name := "config"
	Config = configs.NewConfigmgr()
	Config.AddPath(defines.ConfigPath)

	if err = Config.ReadFile(defines.ConfigFile, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	fmt.Printf("%v initialize\n", name)
	return nil
}

var Config *configs.Configmgr // 配置管理器
