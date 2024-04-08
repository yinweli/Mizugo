package features

import (
	"github.com/yinweli/Mizugo/mizugos/nets"
)

// InitializeNet 初始化網路管理器
func InitializeNet() (err error) {
	name := "net"
	Net = nets.NewNetmgr()
	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

// FinalizeNet 結束網路管理器
func FinalizeNet() {
	if Net != nil {
		Net.Stop()
	} // if
}

var Net *nets.Netmgr // 網路管理器
