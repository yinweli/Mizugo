package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/features"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewKey 建立Key模組
func NewKey() *Key {
	return &Key{
		Module: entitys.NewModule(defines.ModuleIDKey),
		name:   "module key",
		key:    utils.RandDesKeyString(),
	}
}

// Key Key模組
type Key struct {
	*entitys.Module
	name  string // 模組名稱
	key   string // 金鑰
	subID string // 訂閱索引
}

// Awake 喚醒事件
func (this *Key) Awake() error {
	this.subID = this.Entity().Subscribe(entitys.EventSend, this.eventSend)
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyQ), this.procMKeyQ)
	return nil
}

// eventSend 傳送訊息事件
func (this *Key) eventSend(_ any) {
	process, err := utils.CastPointer[procs.Stack](this.Entity().GetProcess())

	if err != nil {
		mizugos.Error(this.name).Message("eventSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().Unsubscribe(this.subID)
}

// procMKeyQ 處理要求Key
func (this *Key) procMKeyQ(message any) {
	rec := features.Ping.Rec()
	defer rec()

	context, ok := message.(*procs.StackContext)

	if ok == false {
		mizugos.Error(this.name).Message("procMKeyQ").EndError(fmt.Errorf("invalid context"))
		return
	} // if

	_, msg, err := procs.StackUnmarshal[messages.MKeyQ](context)

	if err != nil {
		mizugos.Error(this.name).Message("procMKeyQ").EndError(err)
		return
	} // if

	this.sendMKeyA(context, msg, this.key)
	mizugos.Info(this.name).Message("procMKeyQ").KV("key", this.key).End()
}

// sendMKeyA 傳送回應Key
func (this *Key) sendMKeyA(context *procs.StackContext, from *messages.MKeyQ, key string) {
	if err := context.AddRespond(procs.MessageID(messages.MsgID_KeyA), &messages.MKeyA{
		From: from,
		Key:  key,
	}); err != nil {
		mizugos.Error(this.name).Message("sendMKeyA").EndError(err)
	} // if
}
