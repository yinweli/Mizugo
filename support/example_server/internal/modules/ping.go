package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/commons"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewPing 建立Ping模組
func NewPing(pingIncr PingIncr) *Ping {
	return &Ping{
		Module:   entitys.NewModule(defines.ModuleIDPing),
		name:     "module ping(server)",
		pingIncr: pingIncr,
	}
}

// Ping Ping模組
type Ping struct {
	*entitys.Module          // 模組資料
	name            string   // 模組名稱
	pingIncr        PingIncr // 封包計數函式
}

// PingIncr 封包計數函式類型
type PingIncr func() int64

// Awake awake事件
func (this *Ping) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingReq), this.procMsgPingReq)
	return nil
}

// procMsgPingReq 處理要求Ping
func (this *Ping) procMsgPingReq(message any) {
	rec := commons.Ping.Rec()
	defer rec()

	_, msg, err := procs.ProtoDesUnmarshal[*messages.MsgPingReq](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgPingReq").EndError(err)
		return
	} // if

	count := this.pingIncr()
	this.sendMsgPingRes(msg, count)
	mizugos.Info(this.name).Message("procMsgPingReq receive").KV("count", count).End()
}

// sendMsgPingRes 傳送回應Ping
func (this *Ping) sendMsgPingRes(from *messages.MsgPingReq, count int64) {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_PingRes), &messages.MsgPingRes{
		From:  from,
		Count: count,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgPingRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
