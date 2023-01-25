package modules

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/test-client/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client/internal/features"
	"github.com/yinweli/Mizugo/support/test-client/internal/messages"
)

// NewStack 建立Stack模組
func NewStack(disconnect bool, delayTime time.Duration) *Stack {
	return &Stack{
		Module:     entitys.NewModule(defines.ModuleIDStack),
		name:       "module stack",
		disconnect: disconnect,
		delayTime:  delayTime,
	}
}

// Stack Stack模組
type Stack struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	disconnect      bool          // 斷線旗標
	delayTime       time.Duration // 延遲時間
}

// Awake 喚醒事件
func (this *Stack) Awake() error {
	this.Entity().Subscribe(defines.EventKey, this.eventKey)
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyA), this.procMKeyA)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_StackA), this.procMStackA)
	return nil
}

// Start 啟動事件
func (this *Stack) Start() error {
	this.Entity().PublishOnce(defines.EventKey, nil)
	return nil
}

// eventKey key事件
func (this *Stack) eventKey(_ any) {
	this.sendMKeyQ()
}

// event 開始事件
func (this *Stack) eventBegin(_ any) {
	this.sendMStackQ()
}

// procMKeyA 處理回應Key
func (this *Stack) procMKeyA(message any) {
	rec := features.Key.Rec()
	defer rec()

	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMKeyA").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MKeyA](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMKeyA").EndError(err)
		return
	} // if

	process, err := utils.CastPointer[procs.Stack](this.Entity().GetProcess())

	if err != nil {
		mizugos.Error(this.name).Message("procMKeyA").EndError(err)
		return
	} // if

	process.Key([]byte(msg.Key))
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delayTime)
	mizugos.Info(this.name).Message("procMKeyA").KV("key", msg.Key).End()
}

// sendMKeyQ 傳送要求Key
func (this *Stack) sendMKeyQ() {
	context := &procs.StackContext{}

	if err := context.AddRespond(procs.MessageID(messages.MsgID_KeyQ), &messages.MKeyQ{}); err != nil {
		mizugos.Error(this.name).Message("sendMKeyQ").EndError(err)
		return
	} // if

	message, err := procs.StackMarshal(context)

	if err != nil {
		mizugos.Error(this.name).Message("sendMKeyQ").EndError(err)
		return
	} // if

	this.Entity().Send(message)
}

// procMStackA 處理回應Stack
func (this *Stack) procMStackA(message any) {
	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMStackA").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MStackA](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMStackA").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Stack.Add(duration)

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMStackQ()
	} // if

	mizugos.Info(this.name).Message("procMStackA").
		KV("duration", duration).
		KV("count", msg.Count).
		End()
}

// sendMStackQ 傳送要求Stack
func (this *Stack) sendMStackQ() {
	context := &procs.StackContext{}

	if err := context.AddRespond(procs.MessageID(messages.MsgID_StackQ), &messages.MStackQ{
		Time: time.Now().UnixNano(),
	}); err != nil {
		mizugos.Error(this.name).Message("sendMStackQ").EndError(err)
		return
	} // if

	msg, err := procs.StackMarshal(context)

	if err != nil {
		mizugos.Error(this.name).Message("sendMStackQ").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
