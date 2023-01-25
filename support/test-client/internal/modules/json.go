package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client/internal/features"
	"github.com/yinweli/Mizugo/support/test-client/internal/messages"
)

// NewJson 建立Json模組
func NewJson(disconnect bool, delayTime time.Duration) *Json {
	return &Json{
		Module:     entitys.NewModule(defines.ModuleIDJson),
		name:       "module json",
		disconnect: disconnect,
		delayTime:  delayTime,
	}
}

// Json Json模組
type Json struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	disconnect      bool          // 斷線旗標
	delayTime       time.Duration // 延遲時間
}

// Awake 喚醒事件
func (this *Json) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_JsonA), this.procMJsonA)
	return nil
}

// Start 啟動事件
func (this *Json) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delayTime)
	return nil
}

// event 開始事件
func (this *Json) eventBegin(_ any) {
	this.sendMJsonQ()
}

// procMJsonA 處理回應Json
func (this *Json) procMJsonA(message any) {
	_, msg, err := procs.JsonUnmarshal[messages.MJsonA](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMJsonA").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Json.Add(duration)

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMJsonQ()
	} // if

	mizugos.Info(this.name).Message("procMJsonA").
		KV("duration", duration).
		KV("count", msg.Count).
		End()
}

// sendMJsonQ 傳送要求Json
func (this *Json) sendMJsonQ() {
	msg, err := procs.JsonMarshal(procs.MessageID(messages.MsgID_JsonQ), &messages.MJsonQ{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMJsonQ").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
