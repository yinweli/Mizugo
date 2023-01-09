package modules

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_server/internal/commons"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/messages"
)

// NewKey 建立密鑰模組
func NewKey() *Key {
	return &Key{
		Module: entitys.NewModule(defines.ModuleIDKey),
		name:   "module key(server)",
	}
}

// Key 密鑰模組
type Key struct {
	*entitys.Module              // 模組資料
	name            string       // 模組名稱
	event           events.Index // 事件編號
	key             string       // 密鑰
}

// Start start事件
func (this *Key) Start() error {
	index, err := this.Entity().SubEvent(entitys.EventAfterSend, this.afterSend)

	if err != nil {
		return fmt.Errorf("%v start: %w", this.name, err)
	} // if

	this.event = index
	this.key = utils.RandDesKeyString()
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyReq), this.procKeyReq)
	return nil
}

// afterSend afterSend事件
func (this *Key) afterSend(_ any) {
	process, err := utils.CastPointer[procs.ProtoDes](this.Entity().GetProcess())

	if err != nil {
		_ = mizugos.Error(this.name).Message("afterSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().UnsubEvent(this.event)
}

// procPingReq 處理要求密鑰
func (this *Key) procKeyReq(message any) {
	rec := commons.Key.Rec()
	defer rec()

	_, _, err := procs.ProtoDesUnmarshal[*messages.KeyReq](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procKeyReq").EndError(err)
		return
	} // if

	this.sendKeyRes()
	mizugos.Info(this.name).Message("procKeyReq receive").KV("key", this.key).End()
}

// sendKeyRes 傳送回應密鑰
func (this *Key) sendKeyRes() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_KeyRes), &messages.KeyRes{
		Key: this.key,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendKeyRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
