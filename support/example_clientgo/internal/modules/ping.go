package modules

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/commons"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewPing 建立Ping模組
func NewPing(disconnect bool) *Ping {
	return &Ping{
		Module:     entitys.NewModule(defines.ModuleIDPing),
		name:       "module ping(client)",
		disconnect: disconnect,
	}
}

// Ping Ping模組
type Ping struct {
	*entitys.Module              // 模組資料
	name            string       // 模組名稱
	disconnect      bool         // 斷線旗標
	event           events.Index // 事件編號
}

// Awake awake事件
func (this *Ping) Awake() error {
	var err error

	if this.event, err = this.Entity().SubEvent(defines.EventCompleteKey, this.eventCompleteKey); err != nil {
		return fmt.Errorf("%v awake: %w", this.name, err)
	} // if

	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingRes), this.procMsgPingRes)
	return nil
}

// eventCompleteKey completeKey事件
func (this *Ping) eventCompleteKey(_ any) {
	this.sendMsgPingReq()
}

// procMsgPingRes 處理回應Ping
func (this *Ping) procMsgPingRes(message any) {
	_, msg, err := procs.ProtoDesUnmarshal[*messages.MsgPingRes](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgPingRes").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	commons.Ping.Add(duration)
	mizugos.Info(this.name).Message("procMsgPingRes").
		KV("count", msg.Count).
		KV("duration", duration).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} // if
}

// sendMsgPingReq 傳送要求Ping
func (this *Ping) sendMsgPingReq() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_PingReq), &messages.MsgPingReq{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgPingReq").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
