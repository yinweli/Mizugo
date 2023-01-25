package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/messages"
)

// NewStack 建立Stack模組
func NewStack(incr func() int64) *Stack {
	return &Stack{
		Module: entitys.NewModule(defines.ModuleIDStack),
		name:   "module stack",
		incr:   incr,
		key:    utils.RandDesKeyString(),
	}
}

// Stack Stack模組
type Stack struct {
	*entitys.Module
	name  string       // 模組名稱
	incr  func() int64 // 計數函式
	key   string       // 金鑰
	subID string       // 訂閱索引
}

// Awake 喚醒事件
func (this *Stack) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyQ), this.procMKeyQ)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_StackQ), this.procMStackQ)
	return nil
}

// eventSend 傳送訊息事件
func (this *Stack) eventSend(_ any) {
	process, err := utils.CastPointer[procs.Stack](this.Entity().GetProcess())

	if err != nil {
		mizugos.Error(this.name).Message("eventSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().Unsubscribe(this.subID)
}

// procMKeyQ 處理要求Key
func (this *Stack) procMKeyQ(message any) {
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
func (this *Stack) sendMKeyA(context *procs.StackContext, from *messages.MKeyQ, key string) {
	if err := context.AddRespond(procs.MessageID(messages.MsgID_KeyA), &messages.MKeyA{
		From: from,
		Key:  key,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMKeyA").EndError(err)
	} // if
}

// procMStackQ 處理要求Stack
func (this *Stack) procMStackQ(message any) {
	rec := features.Stack.Rec()
	defer rec()

	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMStackQ").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MStackQ](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMStackQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMStackA(context, msg, count)
	mizugos.Info(this.name).Message("procMStackQ").KV("count", count).End()
}

// sendMStackA 傳送回應Stack
func (this *Stack) sendMStackA(context *procs.StackContext, from *messages.MStackQ, count int64) {
	if err := context.AddRespond(procs.MessageID(messages.MsgID_StackA), &messages.MStackA{
		From:  from,
		Count: count,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMStackA").EndError(err)
	} // if
}
