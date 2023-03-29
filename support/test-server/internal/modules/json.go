package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/msgs"
)

const nameJson = "module-json" // 模組名稱

// NewJson 建立Json模組
func NewJson(incr func() int64) *Json {
	return &Json{
		Module: entitys.NewModule(defines.ModuleIDJson),
		incr:   incr,
	}
}

// Json Json模組
type Json struct {
	*entitys.Module
	incr func() int64 // 計數函式
}

// Awake 喚醒處理
func (this *Json) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_JsonQ), this.procMJsonQ)
	return nil
}

// procMJsonQ 處理要求Json
func (this *Json) procMJsonQ(message any) {
	rec := features.Json.Rec()
	defer rec()

	_, msg, err := procs.JsonUnmarshal[msgs.MJsonQ](message)

	if err != nil {
		this.sendMJsonA(nil, msgs.ErrID_JsonUnmarshal, 0)
		features.System.Warn(nameJson).Caller(0).EndError(fmt.Errorf("json procMJsonQ: %w", err))
		return
	} // if

	count := this.incr()
	this.sendMJsonA(msg, msgs.ErrID_Success, count)
	features.System.Info(nameJson).Caller(0).KV("count", count).End()
}

// sendMJsonA 傳送回應Json
func (this *Json) sendMJsonA(from *msgs.MJsonQ, errID msgs.ErrID, count int64) {
	msg, err := procs.JsonMarshal(procs.MessageID(msgs.MsgID_JsonA), &msgs.MJsonA{
		From:  from,
		ErrID: int(errID),
		Count: count,
	})

	if err != nil {
		features.System.Warn(nameJson).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
