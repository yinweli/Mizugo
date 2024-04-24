package modules

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/patterns"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/msgs"
)

// NewProtoRaven 建立ProtoRaven模組
func NewProtoRaven(delay time.Duration, disconnect bool) *ProtoRaven {
	return &ProtoRaven{
		Module:     entitys.NewModule(defines.ModuleIDProtoRaven),
		name:       "protoRaven",
		delay:      delay,
		disconnect: disconnect,
	}
}

// ProtoRaven ProtoRaven模組
type ProtoRaven struct {
	*entitys.Module               // 模組資料
	name            string        // 系統名稱
	delay           time.Duration // 延遲時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒處理
func (this *ProtoRaven) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_ProtoRavenA), this.procMProtoRavenA)
	return nil
}

// Start 啟動處理
func (this *ProtoRaven) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delay)
	return nil
}

// event 開始事件
func (this *ProtoRaven) eventBegin(_ any) {
	this.sendMProtoRavenQ()
}

// procMProtoRavenA 處理回應ProtoRaven
func (this *ProtoRaven) procMProtoRavenA(message any) {
	raven, err := patterns.RavenAParser[msgs.HProtoRaven, msgs.MProtoRavenQ](message)

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("protoRaven procMProtoRavenA: %w", err)).EndFlush()
		return
	} // if

	if raven.ErrID != int32(msgs.ErrID_Success) {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("protoRaven procMProtoRavenA: %v", raven.ErrID)).EndFlush()
		return
	} // if

	respond, ok := raven.GetRespond(&msgs.MProtoRavenA{}).(*msgs.MProtoRavenA)

	if ok == false {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("protoRaven procMProtoRavenA: respond failed")).EndFlush()
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - raven.Request.Time)
	features.MeterProtoRaven.Add(duration)
	features.LogSystem.Get().Info(this.name).KV("count", respond.Count).KV("duration", duration).Caller(0).EndFlush()

	if this.disconnect {
		this.Entity().Stop()
	} else {
		this.sendMProtoRavenQ()
	} // if
}

// sendMProtoRavenQ 傳送要求ProtoRaven
func (this *ProtoRaven) sendMProtoRavenQ() {
	raven, err := patterns.RavenQBuilder(procs.MessageID(msgs.MsgID_ProtoRavenQ),
		&msgs.HProtoRaven{
			Token: "protoRaven",
		},
		&msgs.MProtoRavenQ{
			Time: time.Now().UnixNano(),
		})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(raven)
}
