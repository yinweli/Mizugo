package features

import (
	"github.com/yinweli/Mizugo/mizugo/labels"
)

// InitializeLabel 初始化標籤管理器
func InitializeLabel() (err error) {
	name := "label"
	Label = labels.NewLabelmgr()
	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

var Label *labels.Labelmgr // 標籤管理器
