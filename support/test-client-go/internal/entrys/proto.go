package entrys //nolint:dupl

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/miscs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/modules"
)

// NewProto 建立Proto入口
func NewProto() *Proto {
	return &Proto{
		name: "proto",
	}
}

// Proto Proto入口
type Proto struct {
	name   string      // 入口名稱
	config ProtoConfig // 配置資料
}

// ProtoConfig 配置資料
type ProtoConfig struct {
	Enable     bool          `yaml:"enable"`     // 啟用旗標
	IP         string        `yaml:"ip"`         // 位址
	Port       string        `yaml:"port"`       // 埠號
	Timeout    time.Duration `yaml:"timeout"`    // 超期時間
	Interval   time.Duration `yaml:"interval"`   // 間隔時間
	Count      int           `yaml:"count"`      // 總連線數
	Batch      int           `yaml:"batch"`      // 批次連線數
	Delay      time.Duration `yaml:"delay"`      // 延遲時間
	Disconnect bool          `yaml:"disconnect"` // 斷線旗標
}

// Initialize 初始化處理
func (this *Proto) Initialize() error {
	mizugos.Info(this.name).Caller(0).Message("entry initialize").End()

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if this.config.Enable {
		miscs.GenerateConnection(this.config.Interval, this.config.Count, this.config.Batch, func() {
			mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)
		})
	} // if

	mizugos.Info(this.name).Caller(0).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Proto) Finalize() {
	mizugos.Info(this.name).Caller(0).Message("entry finalize").End()
}

// bind 綁定處理
func (this *Proto) bind(session nets.Sessioner) *nets.Bundle {
	mizugos.Info(this.name).Caller(0).Message("bind").End()
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

	if err := entity.SetProcess(procs.NewProto()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewProto(this.config.Delay, this.config.Disconnect)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.wrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "proto")
	features.Connect.Add(1)
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	mizugos.Error(this.name).Caller(0).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Proto) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
		features.Connect.Add(-1)
	} // if
}

// wrong 錯誤處理
func (this *Proto) wrong(fail bool, err error) {
	if fail {
		mizugos.Error(this.name).Caller(1).EndError(err)
	} else {
		mizugos.Warn(this.name).Caller(1).EndError(err)
	} // if
}
