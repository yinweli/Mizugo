package entrys

import (
	"fmt"
	"path/filepath"

	"github.com/yinweli/Mizugo/cores/msgs"
	"github.com/yinweli/Mizugo/cores/nets"
	"github.com/yinweli/Mizugo/mizugos"
)

// NewEcho 建立回音入口資料
func NewEcho() *Echo {
	return &Echo{
		name: "echo",
	}
}

// Echo 回音入口資料
type Echo struct {
	name   string     // 入口名稱
	config EchoConfig // 設定資料
}

// EchoConfig 回音設定資料
type EchoConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
}

// Initialize 初始化處理
func (this *Echo) Initialize(configPath string) error {
	mizugos.Info(this.name).
		Message("entry initialize").
		End()

	if err := mizugos.Configmgr().ReadFile(filepath.Join(configPath, this.name+".yaml")); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().GetObject(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Netmgr().AddListen(nets.NewTCPListen(this.config.IP, this.config.Port), this)
	mizugos.Info(this.name).
		Message("entry start").
		KV("ip", this.config.IP).
		KV("port", this.config.Port).
		End()
	return nil
}

// Finalize 結束處理
func (this *Echo) Finalize() {
	mizugos.Info(this.name).
		Message("entry stop").
		End()
}

// Bind 綁定處理
func (this *Echo) Bind(session nets.Sessioner) (content nets.Content, err error) {
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	// TODO: add module

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(msgs.NewStringProc()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.Initialize(); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.EntityTagmgr().Add(entity, this.name)
	content.Unbind = func() {
		_ = entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.EntityTagmgr().Del(entity, this.name)
	}
	content.Encode = nil
	content.Decode = nil
	content.Receive = nil

	return content, nil
}

// Error 錯誤處理
func (this *Echo) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
