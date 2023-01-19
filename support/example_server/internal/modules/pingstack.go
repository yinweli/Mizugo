package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/features"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewPingStack 建立PingStack模組
func NewPingStack(incr func() int64) *PingStack {
	return &PingStack{
		Module: entitys.NewModule(defines.ModuleIDPingStack),
		name:   "module pingstack",
		incr:   incr,
		key:    utils.RandDesKeyString(),
	}
}

// PingStack PingStack模組
type PingStack struct {
	*entitys.Module
	name  string       // 模組名稱
	incr  func() int64 // 計數函式
	key   string       // 金鑰
	subID string       // 訂閱索引
}

// Awake 喚醒事件
func (this *PingStack) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyQ), this.procMKeyQ)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingStackQ), this.procMPingStackQ)
	return nil
}

// eventSend 傳送訊息事件
func (this *PingStack) eventSend(_ any) {
	process, err := utils.CastPointer[procs.Stack](this.Entity().GetProcess())

	if err != nil {
		mizugos.Error(this.name).Message("eventSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().Unsubscribe(this.subID)
}

// procMKeyQ 處理要求Key
func (this *PingStack) procMKeyQ(message any) {
	rec := features.Key.Rec()
	defer rec()

	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMKeyQ").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MKeyQ](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMKeyQ").EndError(err)
		return
	} // if

	this.subID = this.Entity().Subscribe(entitys.EventSend, this.eventSend)
	this.sendMKeyA(context, msg, this.key)
	mizugos.Info(this.name).Message("procMKeyQ").KV("key", this.key).End()
}

// sendMKeyA 傳送回應Key
func (this *PingStack) sendMKeyA(context *procs.StackContext, from *messages.MKeyQ, key string) {
	if err := context.AddRespond(procs.MessageID(messages.MsgID_KeyA), &messages.MKeyA{
		From: from,
		Key:  key,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMKeyA").EndError(err)
	} // if
}

// procMPingStackQ 處理要求PingStack
func (this *PingStack) procMPingStackQ(message any) {
	rec := features.Ping.Rec()
	defer rec()

	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMPingStackQ").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MPingStackQ](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMPingStackQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMPingStackA(context, msg, count)
	mizugos.Info(this.name).Message("procMPingStackQ").KV("count", count).End()
}

// sendMPingStackA 傳送回應PingStack
func (this *PingStack) sendMPingStackA(context *procs.StackContext, from *messages.MPingStackQ, count int64) {
	if err := context.AddRespond(procs.MessageID(messages.MsgID_PingStackA), &messages.MPingStackA{
		From:  from,
		Count: count,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMPingStackA").EndError(err)
	} // if
}
