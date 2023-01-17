package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewPing 建立Ping模組
func NewPing(waitKey, waitPing time.Duration, disconnect bool) *Ping {
	return &Ping{
		Module:     entitys.NewModule(defines.ModuleIDPing),
		name:       "module ping(client)",
		waitKey:    waitKey,
		waitPing:   waitPing,
		disconnect: disconnect,
	}
}

// Ping Ping模組
type Ping struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	waitKey         time.Duration // 等待要求Key時間
	waitPing        time.Duration // 等待要求Ping時間
	disconnect      bool          // 斷線旗標
}

// Awake awake事件
func (this *Ping) Awake() error {
	this.Entity().Subscribe(defines.EventPing, this.eventPing)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyRes), this.procMsgKeyRes)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingRes), this.procMsgPingRes)
	return nil
}

// Start start事件
func (this *Ping) Start() error {
	time.Sleep(this.waitKey)
	this.sendMsgKeyReq()
	return nil
}

// eventPing ping事件
func (this *Ping) eventPing(_ any) {
	time.Sleep(this.waitPing)
	this.sendMsgPingReq()
}

// procMsgKeyRes 處理回應密鑰
func (this *Ping) procMsgKeyRes(message any) {
	rec := features.Key.Rec()
	defer rec()

	_, msg, err := procs.ProtoDesUnmarshal[*messages.MsgKeyRes](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgKeyRes").EndError(err)
		return
	} // if

	process, err := utils.CastPointer[procs.ProtoDes](this.Entity().GetProcess())

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgKeyRes").EndError(err)
		return
	} // if

	process.Key([]byte(msg.Key))
	this.Entity().PublishOnce(defines.EventPing, nil)
	mizugos.Info(this.name).Message("procMsgKeyRes").KV("key", msg.Key).End()
}

// sendMsgKeyReq 傳送要求密鑰
func (this *Ping) sendMsgKeyReq() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_KeyReq), &messages.MsgKeyReq{})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgKeyReq").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}

// procMsgPingRes 處理回應Ping
func (this *Ping) procMsgPingRes(message any) {
	_, msg, err := procs.ProtoDesUnmarshal[*messages.MsgPingRes](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgPingRes").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Ping.Add(duration)
	mizugos.Info(this.name).Message("procMsgPingRes").
		KV("count", msg.Count).
		KV("duration", duration).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMsgPingReq()
	} // if
}

// sendMsgPingReq 傳送要求Ping
func (this *Ping) sendMsgPingReq() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_PingReq), &messages.MsgPingReq{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgPingReq").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
