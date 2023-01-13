package entrys

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
)

// NewPing 建立Ping入口資料
func NewPing() *Ping {
	return &Ping{
		name: "pingc",
	}
}

// Ping Ping入口資料
type Ping struct {
	name    string             // 入口名稱
	config  PingConfig         // 設定資料
	connect atomic.Bool        // 連接旗標
	ctx     context.Context    // ctx物件
	cancel  context.CancelFunc // 取消物件
}

// PingConfig 設定資料
type PingConfig struct {
	IP            string        `yaml:"ip"`            // 位址
	Port          string        `yaml:"port"`          // 埠號
	Timeout       time.Duration `yaml:"timeout"`       // 逾期時間(秒)
	Disconnect    bool          `yaml:"disconnect"`    // 斷線旗標
	Reconnect     bool          `yaml:"reconnect"`     // 重連旗標
	ReconnectTime time.Duration `yaml:"reconnectTime"` // 重連檢查時間
	Key           string        `yaml:"key"`           // 密鑰
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

	this.ctx, this.cancel = context.WithCancel(contexts.Ctx())

	go func() {
		timeout := time.NewTicker(this.config.ReconnectTime)

		for {
			select {
			case <-timeout.C:
				if this.connect.Load() {
					continue
				} // if

				if this.config.Reconnect == false {
					continue
				} // if

				this.connect.Store(true)
				mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)

			case <-this.ctx.Done():
				timeout.Stop()
				return
			} // select
		} // for
	}()

	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Ping) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	this.cancel()
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

	if err := entity.AddModule(modules.NewPing(this.config.Disconnect)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.wrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "label echo")
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	this.connect.Store(false)
	session.Stop()
	_ = mizugos.Error(this.name).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Ping) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		this.connect.Store(false)
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *Ping) wrong(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
