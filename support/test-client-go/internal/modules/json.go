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

const nameJson = "module-json" // 模組名稱

// NewJson 建立Json模組
func NewJson(delay time.Duration, disconnect bool) *Json {
	return &Json{
		Module:     entitys.NewModule(defines.ModuleIDJson),
		delay:      delay,
		disconnect: disconnect,
	}
}

// Json Json模組
type Json struct {
	*entitys.Module               // 模組資料
	delay           time.Duration // 延遲時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒處理
func (this *Json) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_JsonA), this.procMJsonA)
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
		features.System.Warn(nameAuth).Caller(0).EndError(fmt.Errorf("json procMJsonA: %w", err))
		return
	} // if

	if msgs.ErrID(msg.ErrID) != msgs.ErrID_Success {
		features.System.Warn(nameAuth).Caller(0).EndError(fmt.Errorf("json procMJsonA: %v", msg.ErrID))
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Json.Add(duration)

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMJsonQ()
	} // if

	features.System.Info(nameJson).Caller(0).KV("duration", duration).KV("count", msg.Count).End()
}

// sendMJsonQ 傳送要求Json
func (this *Json) sendMJsonQ() {
	msg, err := procs.JsonMarshal(procs.MessageID(msgs.MsgID_JsonQ), &msgs.MJsonQ{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		features.System.Warn(nameJson).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
