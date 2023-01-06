package modules

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/procs"
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

	// TODO: 產生密鑰
	// TODO: 傳送密鑰給客戶端
}

// AfterSend 傳送封包後處理函式
func (this *Key) AfterSend() {
	// TODO: 更換處理器, 要怎麼知道封包已送出, 可以改處理器了呢!?
	// TODO: 第一次傳送封包後(也就是傳送了密鑰後)要改成使用正式密鑰的處理器
}

// SendKey 傳送密鑰
func (this *Key) SendKey(key string) error {
	message, err := procs.NewProtoDesMsg(int32(messages.MsgID_ResKey), &messages.MsgResKey{
		Key: key,
	})

	if err != nil {
		// TODO: error process
	} // if

	this.Entity().Send(message)
}
