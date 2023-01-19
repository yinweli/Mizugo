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

// NewPingStack 建立PingStack入口
func NewPingStack() *PingStack {
	return &PingStack{
		name: "pingstack",
	}
}

// PingStack PingStack入口
type PingStack struct {
	name     string          // 入口名稱
	config   PingStackConfig // 配置資料
	listenID nets.ListenID   // 接聽編號
	count    atomic.Int64    // 計數器
}

// PingStackConfig 配置資料
type PingStackConfig struct {
	Enable  bool   `yaml:"enable"`  // 啟用旗標
	IP      string `yaml:"ip"`      // 位址
	Port    string `yaml:"port"`    // 埠號
	KeyInit string `yaml:"keyinit"` // 初始金鑰
}

// Initialize 初始化處理
func (this *PingStack) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if this.config.Enable {
		this.listenID = mizugos.Netmgr().AddListenTCP(this.config.IP, this.config.Port, this.bind, this.unbind, this.wrong)
		mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	} // if

	return nil
}

// Finalize 結束處理
func (this *PingStack) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	mizugos.Netmgr().DelListen(this.listenID)
}

// bind 綁定處理
func (this *PingStack) bind(session nets.Sessioner) *nets.Bundle {
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

	if err := entity.SetProcess(procs.NewStack().Send(entity.Send).Key([]byte(this.config.KeyInit))); err != nil {
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

	if err := entity.AddModule(modules.NewPingStack(this.incr)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.wrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "pingstack")
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	mizugos.Error(this.name).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *PingStack) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *PingStack) wrong(err error) {
	mizugos.Error(this.name).EndError(err)
}

// incr 增加Ping計數
func (this *PingStack) incr() int64 {
	return this.count.Add(1)
}