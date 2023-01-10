package entrys

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

// NewPing 建立Ping入口資料
func NewPing() *Ping {
	return &Ping{
		name: "pingc",
	}
}

// Ping Ping入口資料
type Ping struct {
	name    string      // 入口名稱
	config  PingConfig  // 設定資料
	finish  atomic.Bool // 關閉旗標
	connect atomic.Bool // 連接旗標
}

// PingConfig 設定資料
type PingConfig struct {
	IP            string        // 位址
	Port          string        // 埠號
	Timeout       time.Duration // 逾期時間(秒)
	Disconnect    bool          // 斷線旗標
	Reconnect     bool          // 重連旗標
	ReconnectTime time.Duration // 重連檢查時間
	Key           string        // 密鑰
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
func (this *Ping) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
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

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetProcess(procs.NewProtoDes().Key([]byte(this.config.Key))); err != nil {
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

	if err := entity.Initialize(); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, "label echo")
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	this.connect.Store(false)
	_ = mizugos.Error(this.name).EndError(wrong)
	mizugos.Entitymgr().Del(entity.EntityID())
	session.Stop()
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
