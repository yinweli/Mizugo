package entrys

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/support/example_server/features/defines"
	"github.com/yinweli/Mizugo/support/example_server/features/modules"
)

// NewEchos 建立入口資料
func NewEchos() *Echos {
	return &Echos{
		name: defines.EntryEchos,
	}
}

// Echos 入口資料
type Echos struct {
	name   string        // 入口名稱
	config EchosConfig   // 設定資料
	listen nets.Listener // 接聽物件
}

// EchosConfig 設定資料
type EchosConfig struct {
	IP   string // 位址
	Port string // 埠號
}

// Initialize 初始化處理
func (this *Echos) Initialize() error {
	mizugos.Info(this.name).
		Message("entry initialize").
		End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	this.listen = nets.NewTCPListen(this.config.IP, this.config.Port)
	mizugos.Netmgr().AddListen(this.listen, this)
	mizugos.Info(this.name).
		Message("entry start").
		KV("ip", this.config.IP).
		KV("port", this.config.Port).
		End()
	return nil
}

// Finalize 結束處理
func (this *Echos) Finalize() {
	mizugos.Info(this.name).
		Message("entry stop").
		End()

	if err := this.listen.Stop(); err != nil {
		_ = mizugos.Error(this.name).EndError(err)
	} // if
}

// Bind 綁定處理
func (this *Echos) Bind(session nets.Sessioner) (content nets.Content, err error) {
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(msgs.NewStringProc()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.AddModule(modules.NewEchos()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.Initialize(func() {
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	}); err != nil {
		mizugos.Entitymgr().Del(entity.EntityID())
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEchos)
	content.Unbind = entity.Finalize
	content.Encode = entity.GetProcess().Encode
	content.Decode = entity.GetProcess().Decode
	content.Receive = entity.GetProcess().Process
	return content, nil
}

// Error 錯誤處理
func (this *Echos) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
