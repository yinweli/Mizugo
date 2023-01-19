package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/features"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewPingProto 建立PingProto模組
func NewPingProto(incr func() int64) *PingProto {
	return &PingProto{
		Module: entitys.NewModule(defines.ModuleIDPingProto),
		name:   "module pingproto",
		incr:   incr,
	}
}

// PingProto PingProto模組
type PingProto struct {
	*entitys.Module
	name string       // 模組名稱
	incr func() int64 // 計數函式
}

// Awake 喚醒事件
func (this *PingProto) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingProtoQ), this.procMPingProtoQ)
	return nil
}

// procMPingProtoQ 處理要求PingProto
func (this *PingProto) procMPingProtoQ(message any) {
	rec := features.Ping.Rec()
	defer rec()

	_, msg, err := procs.ProtoUnmarshal[messages.MPingProtoQ](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMPingProtoQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMPingProtoA(msg, count)
	mizugos.Info(this.name).Message("procMPingProtoQ").KV("count", count).End()
}

// sendMPingProtoA 傳送回應PingProto
func (this *PingProto) sendMPingProtoA(from *messages.MPingProtoQ, count int64) {
	msg, err := procs.ProtoMarshal(procs.MessageID(messages.MsgID_PingProtoA), &messages.MPingProtoA{
		From:  from,
		Count: count,
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMPingProtoA").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
