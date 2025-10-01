package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos"
	"github.com/yinweli/Mizugo/v2/support/test-server/internal/defines"
)

// ConfigInitialize 初始化配置
func ConfigInitialize() (err error) {
	mizugos.Config.AddPath(defines.ConfigPath)

	if err = mizugos.Config.ReadFile(defines.ConfigFile, defines.ConfigType); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	if err = mizugos.Config.ReadEnvironment(defines.ConfigEnv); err != nil {
		return fmt.Errorf("initialize: %w", err)
	} // if

	fmt.Println("config initialize")
	return nil
}
