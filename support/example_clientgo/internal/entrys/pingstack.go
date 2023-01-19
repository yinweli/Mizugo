package entrys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/utils"
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
	detector utils.Detector  // 連線檢測器
}

// PingStackConfig 配置資料
type PingStackConfig struct {
	Enable     bool          `yaml:"enable"`     // 啟用旗標
	IP         string        `yaml:"ip"`         // 位址
	Port       string        `yaml:"port"`       // 埠號
	Timeout    time.Duration `yaml:"timeout"`    // 逾期時間(秒)
	Total      int           `yaml:"total"`      // 連線總數
	Batch      int           `yaml:"batch"`      // 連線批數
	Interval   time.Duration `yaml:"interval"`   // 連線間隔時間
	WaitTime   time.Duration `yaml:"waittime"`   // 等待時間
	Disconnect bool          `yaml:"disconnect"` // 斷線旗標
	KeyInit    string        `yaml:"keyinit"`    // 初始金鑰
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
		this.detector.Start(this.config.Total, this.config.Batch, this.config.Interval, func() {
			mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)
		})
	} // if

	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *PingStack) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	this.detector.Stop()
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

	if err := entity.AddModule(modules.NewPingStack(this.config.WaitTime, this.config.Disconnect)); err != nil {
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
