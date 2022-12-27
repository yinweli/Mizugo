package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/support/example_server/features/defines"
)

// NewEchos 建立回音伺服器模組
func NewEchos() *Echos {
	return &Echos{
		Module: entitys.NewModule(1),
		name:   "module echos",
	}
}

// Echos 回音伺服器模組
type Echos struct {
	*entitys.Module
	name string
}

// Start start事件
func (this *Echos) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
}

// ProcMsgEcho 處理回音訊息
func (this *Echos) ProcMsgEcho(messageID msgs.MessageID, message any) {
	_, err := msgs.Cast[msgs.StringMsg](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	this.Entity().Send(message)
}
