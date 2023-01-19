package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewPingProto 建立PingProto模組
func NewPingProto(waitPing time.Duration, disconnect bool) *PingProto {
	return &PingProto{
		Module:     entitys.NewModule(defines.ModuleIDPingProto),
		name:       "module pingproto",
		waitPing:   waitPing,
		disconnect: disconnect,
	}
}

// PingProto PingProto模組
type PingProto struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	waitPing        time.Duration // 等待Ping時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒事件
func (this *PingProto) Awake() error {
	this.Entity().Subscribe(defines.EventPing, this.eventPing)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingProtoA), this.procMPingProtoA)
	return nil
}

// Start 啟動事件
func (this *PingProto) Start() error {
	this.Entity().PublishDelay(defines.EventPing, nil, this.waitPing)
	return nil
}

// eventPing ping事件
func (this *PingProto) eventPing(_ any) {
	this.sendMPingProtoQ()
}

// procMPingProtoA 處理回應PingProto
func (this *PingProto) procMPingProtoA(message any) {
	_, msg, err := procs.ProtoUnmarshal[messages.MPingProtoA](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMPingProtoA").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Ping.Add(duration)
	mizugos.Info(this.name).Message("procMPingProtoA").
		KV("duration", duration).
		KV("count", msg.Count).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMPingProtoQ()
	} // if
}

// sendMPingProtoQ 傳送要求PingProto
func (this *PingProto) sendMPingProtoQ() {
	msg, err := procs.ProtoMarshal(procs.MessageID(messages.MsgID_PingProtoQ), &messages.MPingProtoQ{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMPingProtoQ").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
