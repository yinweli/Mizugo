package modules

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/errs"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/msgs"
)

// NewPList 建立PList模組
func NewPList(delay time.Duration, disconnect bool) *PList {
	return &PList{
		Module:     entitys.NewModule(defines.ModuleIDPList),
		name:       "module plist",
		delay:      delay,
		disconnect: disconnect,
	}
}

// PList PList模組
type PList struct {
	*entitys.Module               // 模組資料
	name            string        // 模組名稱
	delay           time.Duration // 延遲時間
	disconnect      bool          // 斷線旗標
}

// Awake 喚醒處理
func (this *PList) Awake() error {
	this.Entity().Subscribe(defines.EventBegin, this.eventBegin)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_PListA), this.procMPListA)
	return nil
}

// Start 啟動處理
func (this *PList) Start() error {
	this.Entity().PublishDelay(defines.EventBegin, nil, this.delay)
	return nil
}

// event 開始事件
func (this *PList) eventBegin(_ any) {
	this.sendMPListQ()
}

// procMPListA 處理回應PList
func (this *PList) procMPListA(message any) {
	context, ok := message.(*procs.PListContext)

	if ok == false {
		mizugos.Warn(this.name).Caller(0).EndError(errs.Errort(msgs.ErrID_PListContext))
		return
	} // if

	_, msg, err := procs.PListUnmarshal[msgs.MPListA](context)

	if err != nil {
		mizugos.Warn(this.name).Caller(0).EndError(errs.Errore(msgs.ErrID_PListUnmarshal, err))
		return
	} // if

	if msg.ErrID != msgs.ErrID_Success {
		mizugos.Warn(this.name).Caller(0).EndError(errs.Errort(msg.ErrID))
		return
	} // if

	duration := time.Duration(time.Now().UnixNano() - msg.From.Time)
	features.PList.Add(duration)

	if this.disconnect {
		this.Entity().GetSession().Stop()
	} else {
		this.sendMPListQ()
	} // if

	mizugos.Info(this.name).Caller(0).KV("duration", duration).KV("count", msg.Count).End()
}

// sendMPListQ 傳送要求PList
func (this *PList) sendMPListQ() {
	context := &procs.PListContext{}

	if err := context.AddRespond(procs.MessageID(msgs.MsgID_PListQ), &msgs.MPListQ{
		Time: time.Now().UnixNano(),
	}); err != nil {
		mizugos.Warn(this.name).Caller(0).EndError(err)
		return
	} // if

	msg, err := procs.PListMarshal(context)

	if err != nil {
		mizugos.Warn(this.name).Caller(0).EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
