package modules

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/msgs"
)

// NewRaven 建立Raven模組
func NewRaven(delay time.Duration, disconnect bool) *Raven {
	return &Raven{
		Module:     entitys.NewModule(defines.ModuleIDRaven),
		name:       "raven",
		delay:      delay,
		disconnect: disconnect,
	}
}

// Raven Raven模組
type Raven struct {
	*entitys.Module               // 模組資料
	name            string        // 系統名稱
	delay           time.Duration // 延遲時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒處理
func (this *Raven) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(int32(msgs.MsgID_RavenA), this.procMRavenA)
	return nil
}

// Start 啟動處理
func (this *Raven) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delay)
	return nil
}

// event 開始事件
func (this *Raven) eventBegin(_ any) {
	this.sendMRavenQ()
}

// procMRavenA 處理回應Raven
func (this *Raven) procMRavenA(message any) {
	raven, err := procs.RavenCParser[msgs.HRaven, msgs.MRavenQ](message)

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("raven procMRavenA: %w", err)).EndFlush()
		return
	} // if

	if raven.ErrID != int32(msgs.ErrID_Success) {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("raven procMRavenA: %v", raven.ErrID)).EndFlush()
		return
	} // if

	respond := procs.RavenRespondFind(message, &msgs.MRavenA{}).(*msgs.MRavenA)
	duration := time.Duration(time.Now().UnixNano() - raven.Request.Time)
	features.MeterRaven.Add(duration)
	features.LogSystem.Get().Info(this.name).KV("count", respond.Count).KV("duration", duration).Caller(0).EndFlush()

	if this.disconnect {
		this.Entity().Stop()
	} else {
		this.sendMRavenQ()
	} // if
}

// sendMRavenQ 傳送要求Raven
func (this *Raven) sendMRavenQ() {
	raven, err := procs.RavenSBuilder(int32(msgs.MsgID_RavenQ),
		&msgs.HRaven{
			Token: "raven",
		},
		&msgs.MRavenQ{
			Time: time.Now().UnixNano(),
		})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(raven)
}
