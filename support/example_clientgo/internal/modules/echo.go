package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/commons"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewEcho 建立回音模組
func NewEcho(disconnect bool) *Echo {
	return &Echo{
		Module:     entitys.NewModule(1),
		name:       "module echo",
		disconnect: disconnect,
	}
}

// Echo 回音模組
type Echo struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	disconnect      bool   // 斷線旗標
	echo            string // 回音字串
}

// Start start事件
func (this *Echo) Start() {
	this.echo = utils.RandString(defines.EchoCount)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_EchoRes), this.procEchoRes)
	this.sendEchoReq(this.echo)
}

// procEchoRes 處理回應回音
func (this *Echo) procEchoRes(message any) {
	_, msg, err := procs.SimpleUnmarshal[messages.EchoRes](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procEchoRes").EndError(err)
		return
	} // if

	commons.Echo.Add(currentTime() - msg.From.Time)
	mizugos.Info(this.name).Message("procEchoRes").
		KV("equal", this.echo != msg.From.Echo).
		KV("count", msg.Count).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} // if
}

// SendMsgEcho 傳送要求回音
func (this *Echo) sendEchoReq(echo string) {
	msg, err := procs.SimpleMarshal(procs.MessageID(messages.MsgID_EchoReq), &messages.EchoReq{
		Time: currentTime(),
		Echo: echo,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendEchoReq").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

// currentTime 取得現在時間
func currentTime() time.Duration {
	return time.Since(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC))
}
