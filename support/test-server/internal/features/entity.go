package features

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
)

// InitializeEntity 初始化實體管理器
func InitializeEntity() (err error) {
	name := "entity"
	Entity = entitys.NewEntitymgr()
	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

var Entity *entitys.Entitymgr // 實體管理器
