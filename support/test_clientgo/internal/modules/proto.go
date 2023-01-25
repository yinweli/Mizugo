package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/messages"
)

// NewProto 建立Proto模組
func NewProto(disconnect bool, delayTime time.Duration) *Proto {
	return &Proto{
		Module:     entitys.NewModule(defines.ModuleIDProto),
		name:       "module proto",
		disconnect: disconnect,
		delayTime:  delayTime,
	}
}

// Proto Proto模組
type Proto struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	disconnect      bool          // 斷線旗標
	delayTime       time.Duration // 延遲時間
}

// Awake 喚醒事件
func (this *Proto) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_ProtoA), this.procMProtoA)
	return nil
}

// Start 啟動事件
func (this *Proto) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delayTime)
	return nil
}

// event 開始事件
func (this *Proto) eventBegin(_ any) {
	this.sendMProtoQ()
}

// procMProtoA 處理回應Proto
func (this *Proto) procMProtoA(message any) {
	_, msg, err := procs.ProtoUnmarshal[messages.MProtoA](message)

	if err != nil {
		mizugos.Error(this.name).Message("procMProtoA").EndError(err)
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.Proto.Add(duration)

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMProtoQ()
	} // if

	mizugos.Info(this.name).Message("procMProtoA").
		KV("duration", duration).
		KV("count", msg.Count).
		End()
}

// sendMProtoQ 傳送要求Proto
func (this *Proto) sendMProtoQ() {
	msg, err := procs.ProtoMarshal(procs.MessageID(messages.MsgID_ProtoQ), &messages.MProtoQ{
		Time: time.Now().UnixNano(),
	})

	if err != nil {
		mizugos.Error(this.name).Message("sendMProtoQ").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
