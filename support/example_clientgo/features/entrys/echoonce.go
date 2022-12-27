package entrys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/modules"
)

// NewEchoOnce 建立單次回音資料
func NewEchoOnce() *EchoOnce {
	return &EchoOnce{
		name: defines.EntryEchoOnce,
	}
}

// EchoOnce 單次回音資料
type EchoOnce struct {
	name   string         // 入口名稱
	config EchoOnceConfig // 設定資料
}

// EchoOnceConfig 設定資料
type EchoOnceConfig struct {
	IP         string        // 位址
	Port       string        // 埠號
	Timeout    time.Duration // 逾期時間(秒)
	EchoString string        // 回音字串
}

// Initialize 初始化處理
func (this *EchoOnce) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Netmgr().AddConnect(nets.NewTCPConnect(this.config.IP, this.config.Port, this.config.Timeout), this)
	mizugos.Info(this.name).Message("entry start").KV("ip", this.config.IP).KV("port", this.config.Port).KV("timeout", this.config.Timeout).End()
	return nil
}

// Finalize 結束處理
func (this *EchoOnce) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
}

// Bind 綁定處理
func (this *EchoOnce) Bind(session nets.Sessioner) (content nets.Content, err error) {
	mizugos.Info(this.name).Message("session").KV("sessionID", session.SessionID()).End()
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(msgs.NewStringProc()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.AddModule(modules.NewEchoOnce(this.config.EchoString)); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.Initialize(func() {
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	}); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEchoOnce)
	content.Unbind = entity.Finalize
	content.Encode = entity.GetProcess().Encode
	content.Decode = entity.GetProcess().Decode
	content.Receive = entity.GetProcess().Process
	return content, nil
}

// Error 錯誤處理
func (this *EchoOnce) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
