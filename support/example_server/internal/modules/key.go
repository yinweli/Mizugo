package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
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
	*entitys.Module        // 模組資料
	name            string // 模組名稱
}

// Start start事件
func (this *Key) Start() {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyReq), this.procKeyReq)
}

// AfterSend afterSend事件
func (this *Key) AfterSend() {
	// TODO: 做到要利用AfterSend來更換密鑰
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

	key := utils.RandDesKeyString()
	this.sendKeyRes(key)
	mizugos.Info(this.name).Message("procKeyReq receive").KV("key", key).End()
}

// sendKeyRes 傳送回應密鑰
func (this *Key) sendKeyRes(key string) {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_KeyRes), &messages.KeyRes{
		Key: key,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendKeyRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
