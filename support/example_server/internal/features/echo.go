package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/modules"
)

// NewEcho 建立回音入口資料
func NewEcho() *Echo {
	return &Echo{
		name: defines.EntryEcho,
	}
}

// Echo 回音入口資料
type Echo struct {
	name     string        // 入口名稱
	config   EchoConfig    // 設定資料
	listenID nets.ListenID // 接聽編號
}

// EchoConfig 設定資料
type EchoConfig struct {
	IP   string // 位址
	Port string // 埠號
}

// Initialize 初始化處理
func (this *Echo) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	this.listenID = mizugos.Netmgr().AddListenTCP(this.config.IP, this.config.Port, this.bind, this.unbind, this.wrong)
	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Echo) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	mizugos.Netmgr().DelListen(this.listenID)
}

// bind 綁定處理
func (this *Echo) bind(session nets.Sessioner) nets.Bundle {
	mizugos.Info(this.name).Message("bind").End()
	entity := mizugos.Entitymgr().Add()

	var wrong error

	if entity == nil {
		wrong = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetProcess(procs.NewSimple()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewEcho()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEcho)
	session.SetOwner(entity)
	return nets.Bundle{
		Encode:  entity.GetProcess().Encode,
		Decode:  entity.GetProcess().Decode,
		Receive: entity.GetProcess().Process,
	}

Error:
	_ = mizugos.Error(this.name).EndError(wrong)
	mizugos.Entitymgr().Del(entity.EntityID())
	session.Stop()
	return nets.Bundle{}
}

// unbind 解綁處理
func (this *Echo) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *Echo) wrong(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
