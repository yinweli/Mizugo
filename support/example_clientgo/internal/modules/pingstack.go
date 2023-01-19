package modules

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewPingStack 建立PingStack模組
func NewPingStack(waitKey, waitPing time.Duration, disconnect bool) *PingStack {
	return &PingStack{
		Module:     entitys.NewModule(defines.ModuleIDPingStack),
		name:       "module pingstack",
		waitKey:    waitKey,
		waitPing:   waitPing,
		disconnect: disconnect,
	}
}

// PingStack PingStack模組
type PingStack struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	waitKey         time.Duration // 等待Key時間
	waitPing        time.Duration // 等待Ping時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒事件
func (this *PingStack) Awake() error {
	this.Entity().Subscribe(defines.EventKey, this.eventKey)
	this.Entity().Subscribe(defines.EventPing, this.eventPing)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyA), this.procMKeyA)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingStackA), this.procMPingStackA)
	return nil
}

// Start 啟動事件
func (this *PingStack) Start() error {
	this.Entity().PublishDelay(defines.EventKey, nil, this.waitKey)
	return nil
}

// eventKey key事件
func (this *PingStack) eventKey(_ any) {
	this.sendMKeyQ()
}

// eventPing ping事件
func (this *PingStack) eventPing(_ any) {
	this.sendMPingStackQ()
}

// procMKeyA 處理回應Key
func (this *PingStack) procMKeyA(message any) {
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
	this.Entity().PublishDelay(defines.EventPing, nil, this.waitPing)
	mizugos.Info(this.name).Message("procMKeyA").KV("key", msg.Key).End()
}

// sendMKeyQ 傳送要求Key
func (this *PingStack) sendMKeyQ() {
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

// procMPingStackA 處理回應PingStack
func (this *PingStack) procMPingStackA(message any) {
	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMPingStackA").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MPingStackA](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMPingStackA").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Ping.Add(duration)
	mizugos.Info(this.name).Message("procMPingStackA").
		KV("duration", duration).
		KV("count", msg.Count).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMPingStackQ()
	} // if
}

// sendMPingStackQ 傳送要求PingStack
func (this *PingStack) sendMPingStackQ() {
	context := &procs.StackContext{}

	if err := context.AddRespond(procs.MessageID(messages.MsgID_PingStackQ), &messages.MPingStackQ{
		Time: time.Now().UnixNano(),
	}); err != nil {
		mizugos.Error(this.name).Message("sendMPingStackQ").EndError(err)
		return
	} // if

	msg, err := procs.StackMarshal(context)

	if err != nil {
		mizugos.Error(this.name).Message("sendMPingStackQ").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
