package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test_server/internal/defines"
	"github.com/yinweli/Mizugo/support/test_server/internal/features"
	"github.com/yinweli/Mizugo/support/test_server/internal/messages"
)

// NewJson 建立Json模組
func NewJson(incr func() int64) *Json {
	return &Json{
		Module: entitys.NewModule(defines.ModuleIDJson),
		name:   "module json",
		incr:   incr,
	}
}

// Json Json模組
type Json struct {
	*entitys.Module
	name string       // 模組名稱
	incr func() int64 // 計數函式
}

// Awake 喚醒事件
func (this *Json) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_JsonQ), this.procMJsonQ)
	return nil
}

// procMJsonQ 處理要求Json
func (this *Json) procMJsonQ(message any) {
	rec := features.Json.Rec()
	defer rec()

	_, msg, err := procs.JsonUnmarshal[messages.MJsonQ](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMJsonQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMJsonA(msg, count)
	mizugos.Info(this.name).Message("procMJsonQ").KV("count", count).End()
}

// sendMJsonA 傳送回應Json
func (this *Json) sendMJsonA(from *messages.MJsonQ, count int64) {
	msg, err := procs.JsonMarshal(procs.MessageID(messages.MsgID_JsonA), &messages.MJsonA{
		From:  from,
		Count: count,
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMJsonA").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
