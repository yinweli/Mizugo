package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/features"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewPingJson 建立PingJson模組
func NewPingJson(incr func() int64) *PingJson {
	return &PingJson{
		Module: entitys.NewModule(defines.ModuleIDPingJson),
		name:   "module pingjson",
		incr:   incr,
	}
}

// PingJson PingJson模組
type PingJson struct {
	*entitys.Module
	name string       // 模組名稱
	incr func() int64 // 計數函式
}

// Awake 喚醒事件
func (this *PingJson) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingJsonQ), this.procMPingJsonQ)
	return nil
}

// procMPingJsonQ 處理要求PingJson
func (this *PingJson) procMPingJsonQ(message any) {
	rec := features.Ping.Rec()
	defer rec()

	_, msg, err := procs.JsonUnmarshal[messages.MPingJsonQ](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMPingJsonQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMPingJsonA(msg, count)
	mizugos.Info(this.name).Message("procMPingJsonQ").KV("count", count).End()
}

// sendMPingJsonA 傳送回應PingJson
func (this *PingJson) sendMPingJsonA(from *messages.MPingJsonQ, count int64) {
	msg, err := procs.JsonMarshal(procs.MessageID(messages.MsgID_PingJsonA), &messages.MPingJsonA{
		From:  from,
		Count: count,
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMPingJsonA").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
