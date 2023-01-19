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

// NewPingJson 建立PingJson模組
func NewPingJson(waitTime time.Duration, disconnect bool) *PingJson {
	return &PingJson{
		Module:     entitys.NewModule(defines.ModuleIDPingJson),
		name:       "module pingjson",
		waitTime:   waitTime,
		disconnect: disconnect,
	}
}

// PingJson PingJson模組
type PingJson struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	waitTime        time.Duration // 等待時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒事件
func (this *PingJson) Awake() error {
	this.Entity().Subscribe(defines.EventPing, this.eventPing)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_PingJsonA), this.procMPingJsonA)
	return nil
}

// Start 啟動事件
func (this *PingJson) Start() error {
	this.Entity().PublishDelay(defines.EventPing, nil, this.waitTime)
	return nil
}

// eventPing ping事件
func (this *PingJson) eventPing(_ any) {
	this.sendMPingJsonQ()
}

// procMPingJsonA 處理回應PingJson
func (this *PingJson) procMPingJsonA(message any) {
	_, msg, err := procs.JsonUnmarshal[messages.MPingJsonA](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMPingJsonA").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Ping.Add(duration)
	mizugos.Info(this.name).Message("procMPingJsonA").
		KV("duration", duration).
		KV("count", msg.Count).
		End()

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMPingJsonQ()
	} // if
}

// sendMPingJsonQ 傳送要求PingJson
func (this *PingJson) sendMPingJsonQ() {
	msg, err := procs.JsonMarshal(procs.MessageID(messages.MsgID_PingJsonQ), &messages.MPingJsonQ{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMPingJsonQ").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
