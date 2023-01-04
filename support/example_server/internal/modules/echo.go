package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_server/internal/commons"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
)

// NewEcho 建立回音模組
func NewEcho() *Echo {
	return &Echo{
		Module: entitys.NewModule(1),
		name:   "module echo server",
	}
}

// Echo 回音模組
type Echo struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
}

// Start start事件
func (this *Echo) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
}

// ProcMsgEcho 處理回音訊息
func (this *Echo) ProcMsgEcho(messageID procs.MessageID, message any) {
	rec := commons.Echo.Rec()
	defer rec()

	if _, err := utils.CastPointer[procs.StringMsg](message); err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	this.Entity().Send(message)
	mizugos.Info(this.name).Message("ProcMsgEcho receive").End()
}
