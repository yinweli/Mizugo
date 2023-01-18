package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
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
	}
}

// PingStack PingStack模組
type PingStack struct {
	*entitys.Module
	name string       // 模組名稱
	incr func() int64 // 計數函式
}

// Awake 喚醒事件
func (this *PingStack) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingStackQ), this.procMPingStackQ)
	return nil
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

	_, proto, err := context.Unmarshal()

	if err != nil {
		mizugos.Error(this.name).Message("procMPingStackQ").EndError(err)
		return
	} // if

	msg, ok := proto.(*messages.MPingQ)

	if ok == false {
		mizugos.Error(this.name).Message("procMPingStackQ").EndError(fmt.Errorf("invalid message"))
		return
	} // if

	count := this.incr()
	this.sendMPingStackA(context, msg, count)
	mizugos.Info(this.name).Message("procMPingStackQ").KV("count", count).End()
}

// sendMPingStackA 傳送回應PingStack
func (this *PingStack) sendMPingStackA(context *procs.StackContext, from *messages.MPingQ, count int64) {
	if err := context.AddRespond(procs.MessageID(messages.MsgID_PingStackA), &messages.MPingA{
		From:  from,
		Count: count,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMPingStackA").EndError(err)
	} // if
}
