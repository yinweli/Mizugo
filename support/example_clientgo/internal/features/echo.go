package features

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
)

// NewEcho 建立回音入口
func NewEcho() *Echo {
	return &Echo{
		name: "echoc",
	}
}

// Echo 回音入口
type Echo struct {
	name    string      // 入口名稱
	config  EchoConfig  // 設定資料
	finish  atomic.Bool // 關閉旗標
	connect atomic.Bool // 連接旗標
}

// EchoConfig 設定資料
type EchoConfig struct {
	IP            string        // 位址
	Port          string        // 埠號
	Timeout       time.Duration // 逾期時間(秒)
	Disconnect    bool          // 斷線旗標
	Reconnect     bool          // 重連旗標
	ReconnectTime time.Duration // 重連檢查時間
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

	go func() {
		timeout := time.NewTicker(this.config.ReconnectTime)

		for {
			select {
			case <-timeout.C:
				if this.finish.Load() == false {
					if this.connect.Load() {
						continue
					} // if

					if this.config.Reconnect == false {
						continue
					} // if

					this.connect.Store(true)
					mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)
				} // if

			default:
				if this.finish.Load() {
					timeout.Stop()
					return
				} // if
			} // select
		} // for
	}()

	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Echo) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
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

	if err := entity.AddModule(modules.NewEcho(this.config.Disconnect)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "label echo")
	session.SetOwner(entity)
	return nets.Bundle{
		Encode:  entity.GetProcess().Encode,
		Decode:  entity.GetProcess().Decode,
		Receive: entity.GetProcess().Process,
	}

Error:
	this.connect.Store(false)
	_ = mizugos.Error(this.name).EndError(wrong)
	mizugos.Entitymgr().Del(entity.EntityID())
	session.Stop()
	return nets.Bundle{}
}

// unbind 解綁處理
func (this *Echo) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		this.connect.Store(false)
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *Echo) wrong(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
