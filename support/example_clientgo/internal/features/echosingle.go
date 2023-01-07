package features

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
)

// NewEchoSingle 建立單次回音資料
func NewEchoSingle() *EchoSingle {
	return &EchoSingle{
		name: defines.EntryEchoSingle,
	}
}

// EchoSingle 單次回音資料
type EchoSingle struct {
	name   string           // 入口名稱
	config EchoSingleConfig // 設定資料
}

// EchoSingleConfig 設定資料
type EchoSingleConfig struct {
	IP      string        // 位址
	Port    string        // 埠號
	Timeout time.Duration // 逾期時間(秒)
	Message string        // 回音字串
	Repeat  int           // 重複次數
}

// Initialize 初始化處理
func (this *EchoSingle) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.wrong)
	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *EchoSingle) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
}

// bind 綁定處理
func (this *EchoSingle) bind(session nets.Sessioner) nets.Bundle {
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

	if err := entity.AddModule(modules.NewEchoSingle(this.config.Message, this.config.Repeat)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEchoSingle)
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
func (this *EchoSingle) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// wrong 錯誤處理
func (this *EchoSingle) wrong(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
