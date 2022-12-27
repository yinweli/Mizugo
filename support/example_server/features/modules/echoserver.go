package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_server/features/defines"
)

// NewEchoServer 建立回音伺服器模組
func NewEchoServer() *EchoServer {
	return &EchoServer{
		Module: entitys.NewModule(1),
		name:   "module echo server",
	}
}

// EchoServer 回音伺服器模組
type EchoServer struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
}

// Start start事件
func (this *EchoServer) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
}

// ProcMsgEcho 處理回音訊息
func (this *EchoServer) ProcMsgEcho(messageID msgs.MessageID, message any) {
	_, err := utils.CastPointer[msgs.StringMsg](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	this.Entity().Send(message)
	mizugos.Info(this.name).Message("ProcMsgEcho receive").End()
}
