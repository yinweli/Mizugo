package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
)

// NewEchoSingle 建立單次回音模組
func NewEchoSingle(message string, count int) *EchoSingle {
	return &EchoSingle{
		Module:  entitys.NewModule(1),
		name:    "module echo single",
		message: message,
		count:   count,
	}
}

// EchoSingle 單次回音模組
type EchoSingle struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	message         string // 回音字串
	count           int    // 回音次數
}

// Start start事件
func (this *EchoSingle) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
	this.SendMsgEcho()
}

// ProcMsgEcho 處理回音訊息
func (this *EchoSingle) ProcMsgEcho(messageID msgs.MessageID, message any) {
	msg, err := utils.CastPointer[msgs.StringMsg](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	this.count--
	mizugos.Info(this.name).Message("ProcMsgEcho").
		KV("result", msg.Message == this.message).
		KV("count", this.count).
		End()

	if this.count > 0 {
		this.SendMsgEcho()
	} else {
		this.Entity().Finalize()
	} // if
}

// SendMsgEcho 傳送回音訊息
func (this *EchoSingle) SendMsgEcho() {
	this.Entity().Send(&msgs.StringMsg{
		MessageID: defines.MessageIDEcho,
		Message:   this.message,
	})
}
