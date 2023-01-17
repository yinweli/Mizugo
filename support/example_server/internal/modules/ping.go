package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/features"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewPing 建立Ping模組
func NewPing(incr Incr) *Ping {
	return &Ping{
		Module: entitys.NewModule(defines.ModuleIDPing),
		name:   "module ping(server)",
		key:    utils.RandDesKeyString(),
		incr:   incr,
	}
}

// Ping Ping模組
type Ping struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
	key             string // 密鑰
	incr            Incr   // Ping計數函式
	subID           string // 訂閱索引
}

// Incr Ping計數函式類型
type Incr func() int64

// Awake awake事件
func (this *Ping) Awake() error {
	this.subID = this.Entity().Subscribe(entitys.EventSend, this.eventSend)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyReq), this.procMsgKeyReq)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingReq), this.procMsgPingReq)
	return nil
}

// eventSend 傳送事件
func (this *Ping) eventSend(_ any) {
	process, err := utils.CastPointer[procs.ProtoDes](this.Entity().GetProcess())

	if err != nil {
		_ = mizugos.Error(this.name).Message("eventSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().Unsubscribe(this.subID)
}

// procMsgKeyReq 處理要求密鑰
func (this *Ping) procMsgKeyReq(message any) {
	rec := features.Key.Rec()
	defer rec()

	_, _, err := procs.ProtoDesUnmarshal[*messages.MsgKeyReq](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgKeyReq").EndError(err)
		return
	} // if

	this.sendMsgKeyRes()
	mizugos.Info(this.name).Message("procMsgKeyReq").KV("key", this.key).End()
}

// sendMsgKeyRes 傳送回應密鑰
func (this *Ping) sendMsgKeyRes() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_KeyRes), &messages.MsgKeyRes{
		Key: this.key,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgKeyRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

// procMsgPingReq 處理要求Ping
func (this *Ping) procMsgPingReq(message any) {
	rec := features.Ping.Rec()
	defer rec()

	_, msg, err := procs.ProtoDesUnmarshal[*messages.MsgPingReq](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgPingReq").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMsgPingRes(msg, count)
	mizugos.Info(this.name).Message("procMsgPingReq").KV("count", count).End()
}

// sendMsgPingRes 傳送回應Ping
func (this *Ping) sendMsgPingRes(from *messages.MsgPingReq, count int64) {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_PingRes), &messages.MsgPingRes{
		From:  from,
		Count: count,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgPingRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
