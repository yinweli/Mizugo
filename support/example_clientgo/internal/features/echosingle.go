package features

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
)

// NewEchoSingle 建立單次回音資料
func NewEchoSingle() *EchoSingle {
	return &EchoSingle{
		name: defines.EntryEchoSingle,
	}
}

// EchoSingle 單次回音資料
type EchoSingle struct {
	name   string           // 入口名稱
	config EchoSingleConfig // 設定資料
}

// EchoSingleConfig 設定資料
type EchoSingleConfig struct {
	IP      string        // 位址
	Port    string        // 埠號
	Timeout time.Duration // 逾期時間(秒)
	Message string        // 回音字串
	Repeat  int           // 重複次數
}

// Initialize 初始化處理
func (this *EchoSingle) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Netmgr().AddConnect(nets.NewTCPConnect(this.config.IP, this.config.Port, this.config.Timeout), this)
	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *EchoSingle) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
}

// Bind 綁定處理
func (this *EchoSingle) Bind(session nets.Sessioner) (content nets.Content, err error) {
	mizugos.Info(this.name).Message("session").KV("sessionID", session.SessionID()).End()
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(procs.NewSimple()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.AddModule(modules.NewEchoSingle(this.config.Message, this.config.Repeat)); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.Initialize(); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEchoSingle)
	content.Unbind = func() {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	}
	content.Encode = entity.GetProcess().Encode
	content.Decode = entity.GetProcess().Decode
	content.Receive = entity.GetProcess().Process
	return content, nil
}

// Error 錯誤處理
func (this *EchoSingle) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
