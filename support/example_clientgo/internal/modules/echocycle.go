package modules

import (
	"bytes"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
)

// NewEchoCycle 建立循環回音模組
func NewEchoCycle(message string, disconnect bool) *EchoCycle {
	return &EchoCycle{
		Module:     entitys.NewModule(1),
		name:       "module echo cycle",
		message:    []byte(message),
		disconnect: disconnect,
	}
}

// EchoCycle 循環回音模組
type EchoCycle struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	message         []byte // 回音資料
	disconnect      bool   // 斷線旗標
	count           int    // 回音次數
}

// Start start事件
func (this *EchoCycle) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
	this.SendMsgEcho()
}

// ProcMsgEcho 處理回音訊息
func (this *EchoCycle) ProcMsgEcho(messageID procs.MessageID, message any) {
	msg, err := utils.CastPointer[procs.SimpleMsg](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	this.count++
	mizugos.Info(this.name).Message("ProcMsgEcho").
		KV("result", bytes.Equal(msg.Message, this.message)).
		KV("count", this.count).
		End()

	if this.disconnect == false {
		this.SendMsgEcho()
	} else {
		this.Entity().GetSession().Stop()
	} // if
}

// SendMsgEcho 傳送回音訊息
func (this *EchoCycle) SendMsgEcho() {
	this.Entity().Send(&procs.SimpleMsg{
		MessageID: defines.MessageIDEcho,
		Message:   this.message,
	})
}
