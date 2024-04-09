package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/pools"
)

// PoolInitialize 初始化執行緒池
func PoolInitialize() (err error) {
	config := &pools.Config{
		Logger: func(format string, args ...any) {
			LogSystem.Get().Error("pool").Message(format, args...).Caller(1).EndFlush()
		},
	}

	if err = mizugos.Config.Unmarshal("pool", config); err != nil {
		return fmt.Errorf("pool initialize: %w", err)
	} // if

	if err = mizugos.Pool.Initialize(config); err != nil {
		return fmt.Errorf("pool initialize: %w", err)
	} // if

	LogSystem.Get().Info("pool").Message("initialize").EndFlush()
	return nil
}
