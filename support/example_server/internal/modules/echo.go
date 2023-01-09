package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/commons"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewEcho 建立回音模組
func NewEcho(echoIncr EchoIncr) *Echo {
	return &Echo{
		Module:   entitys.NewModule(defines.ModuleIDEcho),
		name:     "module echo(server)",
		echoIncr: echoIncr,
	}
}

// Echo 回音模組
type Echo struct {
	*entitys.Module          // 模組資料
	name            string   // 模組名稱
	echoIncr        EchoIncr // 封包計數函式
}

// EchoIncr 封包計數函式類型
type EchoIncr func() int64

// Start start事件
func (this *Echo) Start() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_EchoReq), this.procMsgEchoReq)
	return nil
}

// procMsgEchoReq 處理要求回音
func (this *Echo) procMsgEchoReq(message any) {
	rec := commons.Echo.Rec()
	defer rec()

	_, msg, err := procs.SimpleUnmarshal[messages.MsgEchoReq](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgEchoReq").EndError(err)
		return
	} // if

	count := this.echoIncr()
	this.sendMsgEchoRes(msg, count)
	mizugos.Info(this.name).Message("procMsgEchoReq receive").KV("count", count).End()
}

// sendMsgEchoRes 傳送回應回音
func (this *Echo) sendMsgEchoRes(from *messages.MsgEchoReq, count int64) {
	msg, err := procs.SimpleMarshal(procs.MessageID(messages.MsgID_EchoRes), &messages.MsgEchoRes{
		From:  *from,
		Count: count,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgEchoRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
