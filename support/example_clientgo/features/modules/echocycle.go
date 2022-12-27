package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/defines"
)

// NewEchoCycle 建立循環回音模組
func NewEchoCycle(echoString string) *EchoCycle {
	return &EchoCycle{
		Module:     entitys.NewModule(1),
		name:       "module echo cycle",
		echoString: echoString,
	}
}

// EchoCycle 循環回音模組
type EchoCycle struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	echoString      string // 回音字串
	echoCount       int    // 回音次數
}

// Start start事件
func (this *EchoCycle) Start() {
	this.Entity().AddMessage(defines.MessageIDEcho, this.ProcMsgEcho)
	this.SendMsgEcho()
}

// ProcMsgEcho 處理回音訊息
func (this *EchoCycle) ProcMsgEcho(messageID msgs.MessageID, message any) {
	msg, err := utils.CastPointer[msgs.StringMsg](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("ProcMsgEcho").EndError(err)
		return
	} // if

	this.echoCount++
	mizugos.Info(this.name).Message("ProcMsgEcho").
		KV("result", msg.Message == this.echoString).
		KV("count", this.echoCount).
		End()
	this.SendMsgEcho()
}

// SendMsgEcho 傳送回音訊息
func (this *EchoCycle) SendMsgEcho() {
	this.Entity().Send(&msgs.StringMsg{
		MessageID: defines.MessageIDEcho,
		Message:   this.echoString,
	})
}
