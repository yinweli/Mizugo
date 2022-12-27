package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/defines"
)

// NewEchoOnce 建立單次回音模組
func NewEchoOnce(echoString string) *EchoOnce {
	return &EchoOnce{
		Module:     entitys.NewModule(1),
		name:       "module echo once",
		echoString: echoString,
	}
}

// EchoOnce 單次回音模組
type EchoOnce struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	echoString      string // 回音字串
}

// Start start事件
func (this *EchoOnce) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
	this.SendMsgEcho()
}

// ProcMsgEcho 處理回音訊息
func (this *EchoOnce) ProcMsgEcho(messageID msgs.MessageID, message any) {
	msg, err := utils.CastPointer[msgs.StringMsg](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	mizugos.Info(this.name).Message("ProcMsgEcho").KV("result", msg.Message == this.echoString).End()
	this.Entity().Finalize()
}

// SendMsgEcho 傳送回音訊息
func (this *EchoOnce) SendMsgEcho() {
	this.Entity().Send(&msgs.StringMsg{
		MessageID: defines.MessageIDEcho,
		Message:   this.echoString,
	})
}
