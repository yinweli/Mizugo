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

// NewJson 建立Json模組
func NewJson(delay time.Duration, disconnect bool) *Json {
	return &Json{
		Module:     entitys.NewModule(defines.ModuleIDJson),
		name:       "json",
		delay:      delay,
		disconnect: disconnect,
	}
}

// Json Json模組
type Json struct {
	*entitys.Module               // 模組資料
	name            string        // 系統名稱
	delay           time.Duration // 延遲時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒處理
func (this *Json) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(int32(msgs.MsgID_JsonA), this.procMJsonA)
	return nil
}

// Start 啟動處理
func (this *Json) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delay)
	return nil
}

// event 開始事件
func (this *Json) eventBegin(_ any) {
	this.sendMJsonQ()
}

// procMJsonA 處理回應Json
func (this *Json) procMJsonA(message any) {
	_, msg, err := procs.JsonUnmarshal[msgs.MJsonA](message)

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("json procMJsonA: %w", err)).EndFlush()
		return
	} // if

	if msgs.ErrID(msg.ErrID) != msgs.ErrID_Success { //nolint:gosec
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(fmt.Errorf("json procMJsonA: %v", msg.ErrID)).EndFlush()
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.MeterJson.Add(duration)
	features.LogSystem.Get().Info(this.name).KV("count", msg.Count).KV("duration", duration).Caller(0).EndFlush()

	if this.disconnect {
		this.Entity().Stop()
	} else {
		this.sendMJsonQ()
	} // if
}

// sendMJsonQ 傳送要求Json
func (this *Json) sendMJsonQ() {
	msg, err := procs.JsonMarshal(int32(msgs.MsgID_JsonQ), &msgs.MJsonQ{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		features.LogSystem.Get().Warn(this.name).Caller(0).Error(err).EndFlush()
		return
	} // if

	this.Entity().Send(msg)
}
