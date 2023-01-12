package modules

import (
	"fmt"

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
	subID           string // 訂閱索引
	key             string // 密鑰
}

// Awake awake事件
func (this *Key) Awake() (err error) {
	if this.subID, err = this.Entity().Subscribe(entitys.EventAfterSend, this.eventAfterSend); err != nil {
		return fmt.Errorf("%v awake: %w", this.name, err)
	} // if

	this.Entity().AddMessage(procs.MessageID(messages.MsgID_KeyReq), this.procMsgKeyReq)
	this.key = utils.RandDesKeyString()
	return nil
}

// eventAfterSend afterSend事件
func (this *Key) eventAfterSend(_ any) {
	process, err := utils.CastPointer[procs.ProtoDes](this.Entity().GetProcess())

	if err != nil {
		_ = mizugos.Error(this.name).Message("eventAfterSend").EndError(err)
		return
	} // if

	process.Key([]byte(this.key))
	this.Entity().Unsubscribe(this.subID)
}

// procMsgKeyReq 處理要求密鑰
func (this *Key) procMsgKeyReq(message any) {
	rec := commons.Key.Rec()
	defer rec()

	_, _, err := procs.ProtoDesUnmarshal[*messages.MsgKeyReq](message)

	if err != nil {
		_ = mizugos.Error(this.name).Message("procMsgKeyReq").EndError(err)
		return
	} // if

	this.sendMsgKeyRes()
	mizugos.Info(this.name).Message("procMsgKeyReq receive").KV("key", this.key).End()
}

// sendMsgKeyRes 傳送回應密鑰
func (this *Key) sendMsgKeyRes() {
	msg, err := procs.ProtoDesMarshal(procs.MessageID(messages.MsgID_KeyRes), &messages.MsgKeyRes{
		Key: this.key,
	})

	if err != nil {
		_ = mizugos.Error(this.name).Message("sendMsgKeyRes").EndError(err)
		return
	} // if

	this.Entity().Send(msg)
}
