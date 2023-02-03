package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/msgs"
)

// NewProto 建立Proto模組
func NewProto(incr func() int64) *Proto {
	return &Proto{
		Module: entitys.NewModule(defines.ModuleIDProto),
		name:   "module proto",
		incr:   incr,
	}
}

// Proto Proto模組
type Proto struct {
	*entitys.Module
	name string       // 模組名稱
	incr func() int64 // 計數函式
}

// Awake 喚醒事件
func (this *Proto) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_ProtoQ), this.procMProtoQ)
	return nil
}

// procMProtoQ 處理要求Proto
func (this *Proto) procMProtoQ(message any) {
	rec := features.Proto.Rec()
	defer rec()

	_, msg, err := procs.ProtoUnmarshal[msgs.MProtoQ](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMProtoQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMProtoA(msg, count)
	mizugos.Info(this.name).Message("procMProtoQ").KV("count", count).End()
}

// sendMProtoA 傳送回應Proto
func (this *Proto) sendMProtoA(from *msgs.MProtoQ, count int64) {
	msg, err := procs.ProtoMarshal(procs.MessageID(msgs.MsgID_ProtoA), &msgs.MProtoA{
		From:  from,
		Count: count,
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMProtoA").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
