package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewEcho 建立回音模組
func NewEcho(disconnect bool) *Echo {
	return &Echo{
		Module:     entitys.NewModule(defines.ModuleIDEcho),
		name:       "module echo(client)",
		disconnect: disconnect,
		echo:       utils.RandString(defines.EchoCount),
	}
}

// Echo 回音模組
type Echo struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	disconnect      bool   // 斷線旗標
	echo            string // 回音字串
}

// Awake awake事件
func (this *Echo) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_EchoRes), this.procMsgEchoRes)
	return nil
}

// Start start事件
func (this *Echo) Start() error {
	this.sendMsgEchoReq(this.echo)
	return nil
}

// procMsgEchoRes 處理回應回音
func (this *Echo) procMsgEchoRes(message any) {
	_, msg, err := procs.SimpleUnmarshal[messages.MsgEchoRes](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgEchoRes").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Echo.Add(duration)
	mizugos.Info(this.name).Message("procMsgEchoRes").
		KV("equal", this.echo == msg.From.Echo).
		KV("count", msg.Count).
		KV("duration", duration).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} // if
}

// sendMsgEchoReq 傳送要求回音
func (this *Echo) sendMsgEchoReq(echo string) {
	msg, err := procs.SimpleMarshal(procs.MessageID(messages.MsgID_EchoReq), &messages.MsgEchoReq{
		Time: time.Now().UnixNano(),
		Echo: echo,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgEchoReq").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
