package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos/entitys"
	"github.com/yinweli/Mizugo/v2/mizugos/procs"
	"github.com/yinweli/Mizugo/v2/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/v2/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/v2/support/test-server/msgs"
)

// NewRaven 建立Raven模組
func NewRaven() *Raven {
	return &Raven{
		Module: entitys.NewModule(defines.ModuleIDRaven),
		name:   "raven",
	}
}

// Raven Raven模組
type Raven struct {
	*entitys.Module
	name string // 系統名稱
}

// Awake 喚醒處理
func (this *Raven) Awake() error {
	this.Entity().AddMessage(int32(msgs.MsgID_RavenQ), this.procMRavenQ)
	return nil
}

// procMRavenQ 處理要求Raven
func (this *Raven) procMRavenQ(message any) {
	raven, err := procs.RavenSParser[*msgs.HRaven, *msgs.MRavenQ](message)

	if err != nil {
		this.sendMRavenA(nil, nil, msgs.ErrID_RavenUnmarshal, 0)
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("raven procMRavenQ: %w", err)).EndFlush()
		return
	} // if

	count := features.RavenCounter.Add(1)
	this.sendMRavenA(raven.Header, raven.Request, msgs.ErrID_Success, count)
	features.LogSystem.Get().Info(this.name).KV("count", count).Caller(0).EndFlush()
}

// sendMRavenA 傳送回應Raven
func (this *Raven) sendMRavenA(header *msgs.HRaven, request *msgs.MRavenQ, errID msgs.ErrID, count int64) {
	raven, err := procs.RavenCBuilder(int32(msgs.MsgID_RavenA), int32(errID), header, request, &msgs.MRavenA{
		Count: count,
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(raven)
}
