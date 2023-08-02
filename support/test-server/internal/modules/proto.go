package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/msgs"
)

// NewProto 建立Proto模組
func NewProto(incr func(int64) int64) *Proto {
	return &Proto{
		Module: entitys.NewModule(defines.ModuleIDProto),
		name:   "proto",
		incr:   incr,
	}
}

// Proto Proto模組
type Proto struct {
	*entitys.Module
	name string            // 系統名稱
	incr func(int64) int64 // 計數函式
}

// Awake 喚醒處理
func (this *Proto) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_ProtoQ), this.procMProtoQ)
	return nil
}

// procMProtoQ 處理要求Proto
func (this *Proto) procMProtoQ(message any) {
	rec := features.MeterProto.Rec()
	defer rec()

	_, msg, err := procs.ProtoUnmarshal[msgs.MProtoQ](message)

	if err != nil {
		this.sendMProtoA(nil, msgs.ErrID_ProtoUnmarshal, 0)
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("proto procMProtoQ: %w", err)).EndFlush()
		return
	} // if

	count := this.incr(1)
	this.sendMProtoA(msg, msgs.ErrID_Success, count)
	features.LogSystem.Get().Info(this.name).KV("count", count).Caller(0).EndFlush()
}

// sendMProtoA 傳送回應Proto
func (this *Proto) sendMProtoA(from *msgs.MProtoQ, errID msgs.ErrID, count int64) {
	msg, err := procs.ProtoMarshal(procs.MessageID(msgs.MsgID_ProtoA), &msgs.MProtoA{
		From:  from,
		ErrID: errID,
		Count: count,
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(msg)
}
