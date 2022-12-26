package module

import (
	"fmt"

	"github.com/yinweli/Mizugo/cores/entitys"
)

// NewEchos 建立回音伺服器模組
func NewEchos() *Echos {
	return &Echos{
		Module: entitys.NewModule(1),
	}
}

// Echos 回音伺服器模組
type Echos struct {
	*entitys.Module
}

// Start start事件
func (this *Echos) Start() {
	fmt.Println("echos module start")
}
