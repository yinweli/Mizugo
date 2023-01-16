package entrys

import (
	"fmt"
	"sync/atomic"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/modules"
)

// NewPing 建立Ping入口
func NewPing() *Ping {
	return &Ping{
		name: "pings",
	}
}

// Ping Ping入口
type Ping struct {
	name      string        // 入口名稱
	config    PingConfig    // 配置資料
	listenID  nets.ListenID // 接聽編號
	pingCount atomic.Int64  // 封包計數
}

// PingConfig 配置資料
type PingConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// Initialize 初始化處理
func (this *Ping) Initialize() error {
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
func (this *Ping) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	mizugos.Netmgr().DelListen(this.listenID)
}

// bind 綁定處理
func (this *Ping) bind(session nets.Sessioner) *nets.Bundle {
	mizugos.Info(this.name).Message("bind").End()
	entity := mizugos.Entitymgr().Add()

	var wrong error

	if entity == nil {
		wrong = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err := entity.SetModulemgr(entitys.NewModulemgr()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetEventmgr(events.NewEventmgr(defines.EventCapacity)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetProcess(procs.NewProtoDes().Key([]byte(this.config.Key))); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewKey()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewPing(func() int64 {
		return this.pingCount.Add(1)
	})); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.wrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "label ping")
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	_ = mizugos.Error(this.name).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Ping) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *Ping) wrong(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
