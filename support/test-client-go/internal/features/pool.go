package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/pools"
)

// InitializePool 初始化執行緒池管理器
func InitializePool() (err error) {
	name := "pool"
	config := &pools.Config{
		Logger: func(format string, args ...any) {
			LogSystem.Get().Error(name).Message(format, args...).Caller(1).EndFlush()
		},
	}
	Pool = pools.DefaultPool

	if err = Config.Unmarshal(name, config); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	if err = Pool.Initialize(config); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

// FinalizePool 結束執行緒池管理器
func FinalizePool() {
	if Pool != nil {
		Pool.Finalize()
	} // if
}

var Pool *pools.Poolmgr // 執行緒池管理器
