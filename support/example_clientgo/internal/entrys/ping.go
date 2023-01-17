package entrys

import (
	"context"
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
)

// NewPing 建立Ping入口
func NewPing() *Ping {
	return &Ping{
		name: "ping",
	}
}

// Ping Ping入口
type Ping struct {
	name    string     // 入口名稱
	config  PingConfig // 配置資料
	connect connect    // 連線檢測
}

// PingConfig 配置資料
type PingConfig struct {
	IP         string        `yaml:"ip"`         // 位址
	Port       string        `yaml:"port"`       // 埠號
	Timeout    time.Duration `yaml:"timeout"`    // 逾期時間(秒)
	InitKey    string        `yaml:"initkey"`    // 初始密鑰
	Count      int           `yaml:"count"`      // 連線總數
	Interval   time.Duration `yaml:"interval"`   // 連線間隔時間
	WaitKey    time.Duration `yaml:"waitkey"`    // 等待要求Key時間
	WaitPing   time.Duration `yaml:"waitping"`   // 等待要求Ping時間
	Disconnect bool          `yaml:"disconnect"` // 斷線旗標
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

	this.connect.start(this.config.Count, this.config.Interval, func() {
		mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)
	})

	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Ping) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	this.connect.stop()
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

	if err := entity.SetProcess(procs.NewProtoDes().Key([]byte(this.config.InitKey))); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewPing(this.config.WaitKey, this.config.WaitPing, this.config.Disconnect)); err != nil {
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

	this.connect.notice()
	session.Stop()
	_ = mizugos.Error(this.name).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Ping) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		this.connect.notice()
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *Ping) wrong(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}

// connect 連線檢測器
type connect struct {
	notify chan any // 通知通道
	cancel func()   // 取消物件
}

// start 啟動連線檢測
func (this *connect) start(count int, interval time.Duration, done func()) {
	mizugos.Poolmgr().Submit(func() {
		conn := func() {
			session := mizugos.Netmgr().Status().Session
			features.Connect.Set(int64(session))

			if session < count {
				done()
			} // if
		}
		timeout := time.NewTicker(interval)
		ctx, cancel := context.WithCancel(contexts.Ctx())
		this.notify = make(chan any, 1)
		this.cancel = cancel

		for {
			select {
			case <-this.notify:
				conn()

			case <-timeout.C:
				conn()

			case <-ctx.Done():
				timeout.Stop()
				return
			} // select
		} // for
	})
}

// stop 停止連線檢測
func (this *connect) stop() {
	this.cancel()
}

// notice 通知連線變化
func (this *connect) notice() {
	this.notify <- nil
}
