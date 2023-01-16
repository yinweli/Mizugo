package modules

import (
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/messages"
)

// NewKey 建立密鑰模組
func NewKey() *Key {
	return &Key{
		Module: entitys.NewModule(defines.ModuleIDKey),
		name:   "module key(client)",
	}
}

// Key 密鑰模組
type Key struct {
	*entitys.Module        // 模組資料
	name            string // 模組名稱
}

// Awake awake事件
func (this *Key) Awake() error {
	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyRes), this.procMsgKeyRes)
	return nil
}

// Start start事件
func (this *Key) Start() error {
	this.sendMsgKeyReq()
	return nil
}

// procMsgKeyRes 處理回應密鑰
func (this *Key) procMsgKeyRes(message any) {
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
	this.Entity().PublishOnce(defines.EventCompleteKey, nil)
	mizugos.Info(this.name).Message("procMsgKeyRes receive").KV("key", msg.Key).End()
}

// sendMsgKeyReq 傳送要求密鑰
func (this *Key) sendMsgKeyReq() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_KeyReq), &messages.MsgKeyReq{})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgKeyReq").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
