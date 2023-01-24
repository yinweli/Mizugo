package entrys //nolint:dupl

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/miscs"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/modules"
)

// NewProto 建立Proto入口
func NewProto() *Proto {
	return &Proto{
		name: "proto",
	}
}

// Proto Proto入口
type Proto struct {
	name      string           // 入口名稱
	config    ProtoConfig      // 配置資料
	generator *miscs.Generator // 連線產生物件
}

// ProtoConfig 配置資料
type ProtoConfig struct {
	Enable     bool          `yaml:"enable"`     // 啟用旗標
	IP         string        `yaml:"ip"`         // 位址
	Port       string        `yaml:"port"`       // 埠號
	Timeout    time.Duration `yaml:"timeout"`    // 逾期時間(秒)
	Max        int           `yaml:"max"`        // 最大連線數
	Batch      int           `yaml:"batch"`      // 批次連線數
	Baseline   time.Duration `yaml:"baseline"`   // 基準時間
	Interval   time.Duration `yaml:"interval"`   // 間隔時間
	Disconnect bool          `yaml:"disconnect"` // 斷線旗標
	DelayTime  time.Duration `yaml:"delayTime"`  // 延遲時間
}

// Initialize 初始化處理
func (this *Proto) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if this.config.Enable {
		this.generator = miscs.NewGenerator(this.config.Max, this.config.Batch, this.config.Baseline, this.config.Interval, func() {
			mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)
		})
		this.generator.Start()
	} // if

	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Proto) Finalize() {
	this.generator.Stop()
	mizugos.Info(this.name).Message("entry finalize").End()
}

// bind 綁定處理
func (this *Proto) bind(session nets.Sessioner) *nets.Bundle {
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

	if err := entity.SetProcess(procs.NewProto()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewProto(this.config.Disconnect, this.config.DelayTime, this.generator.Report)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.wrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "proto")
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
func (this *Proto) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *Proto) wrong(err error) {
	mizugos.Error(this.name).EndError(err)
}
