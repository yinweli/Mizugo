package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/msgs"
)

// NewPList 建立PList模組
func NewPList(incr func() int64) *PList {
	return &PList{
		Module: entitys.NewModule(defines.ModuleIDPList),
		name:   "module plist",
		incr:   incr,
		key:    cryptos.RandDesKeyString(),
	}
}

// PList PList模組
type PList struct {
	*entitys.Module
	name  string       // 模組名稱
	incr  func() int64 // 計數函式
	key   string       // 金鑰
	subID string       // 訂閱索引
}

// Awake 喚醒事件
func (this *PList) Awake() error {
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_KeyQ), this.procMKeyQ)
	this.Entity().AddMessage(procs.MessageID(msgs.MsgID_PListQ), this.procMPListQ)
	return nil
}

// eventSend 傳送訊息事件
func (this *PList) eventSend(_ any) {
	process, err := utils.CastPointer[procs.PList](this.Entity().GetProcess())

	if err != nil {
		mizugos.Error(this.name).Message("eventSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().Unsubscribe(this.subID)
}

// procMKeyQ 處理要求Key
func (this *PList) procMKeyQ(message any) {
	rec := features.Key.Rec()
	defer rec()

	context, ok := message.(*procs.PListContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMKeyQ").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.PListUnmarshal[msgs.MKeyQ](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMKeyQ").EndError(err)
		return
	} // if

	this.subID = this.Entity().Subscribe(entitys.EventSend, this.eventSend)
	this.sendMKeyA(context, msg, this.key)
	mizugos.Info(this.name).Message("procMKeyQ").KV("key", this.key).End()
}

// sendMKeyA 傳送回應Key
func (this *PList) sendMKeyA(context *procs.PListContext, from *msgs.MKeyQ, key string) {
	if err := context.AddRespond(procs.MessageID(msgs.MsgID_KeyA), &msgs.MKeyA{
		From: from,
		Key:  key,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMKeyA").EndError(err)
	} // if
}

// procMPListQ 處理要求PList
func (this *PList) procMPListQ(message any) {
	rec := features.PList.Rec()
	defer rec()

	context, ok := message.(*procs.PListContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMPListQ").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.PListUnmarshal[msgs.MPListQ](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMPListQ").EndError(err)
		return
	} // if

	count := this.incr()
	this.sendMPListA(context, msg, count)
	mizugos.Info(this.name).Message("procMPListQ").KV("count", count).End()
}

// sendMPListA 傳送回應PList
func (this *PList) sendMPListA(context *procs.PListContext, from *msgs.MPListQ, count int64) {
	if err := context.AddRespond(procs.MessageID(msgs.MsgID_PListA), &msgs.MPListA{
		From:  from,
		Count: count,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMPListA").EndError(err)
	} // if
}
