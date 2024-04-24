package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/patterns"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/msgs"
)

// NewProtoRaven 建立ProtoRaven模組
func NewProtoRaven() *ProtoRaven {
	return &ProtoRaven{
		Module: entitys.NewModule(defines.ModuleIDProtoRaven),
		name:   "protoRaven",
	}
}

// ProtoRaven ProtoRaven模組
type ProtoRaven struct {
	*entitys.Module
	name string // 系統名稱
}

// Awake 喚醒處理
func (this *ProtoRaven) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_ProtoRavenQ), this.procMProtoRavenQ)
	return nil
}

// procMProtoRavenQ 處理要求ProtoRaven
func (this *ProtoRaven) procMProtoRavenQ(message any) {
	rec := features.MeterProtoRaven.Rec()
	defer rec()
	raven, err := patterns.RavenQParser[msgs.HProtoRaven, msgs.MProtoRavenQ](message)

	if err != nil {
		this.sendMProtoRavenA(nil, nil, msgs.ErrID_ProtoRavenUnmarshal, 0)
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("raven procMProtoRavenQ: %w", err)).EndFlush()
		return
	} // if

	count := features.ProtoRavenCounter.Add(1)
	this.sendMProtoRavenA(raven.Header, raven.Request, msgs.ErrID_Success, count)
	features.LogSystem.Get().Info(this.name).KV("count", count).Caller(0).EndFlush()
}

// sendMProtoRavenA 傳送回應ProtoRaven
func (this *ProtoRaven) sendMProtoRavenA(header *msgs.HProtoRaven, request *msgs.MProtoRavenQ, errID msgs.ErrID, count int64) {
	raven, err := patterns.RavenABuilder(procs.MessageID(msgs.MsgID_ProtoRavenA), int32(errID), header, request, &msgs.MProtoRavenA{
		Count: count,
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(raven)
}
