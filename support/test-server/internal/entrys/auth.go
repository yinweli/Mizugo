package entrys

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

const nameAuth = "auth" // 入口名稱

// NewAuth 建立Auth入口
func NewAuth() *Auth {
	return &Auth{}
}

// Auth Auth入口
type Auth struct {
	config   AuthConfig    // 配置資料
	listenID nets.ListenID // 接聽編號
}

// AuthConfig 配置資料
type AuthConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
}

// Initialize 初始化處理
func (this *Auth) Initialize() error {
	mizugos.Info(defines.LogSystem, nameAuth).Caller(0).Message("entry initialize").End()

	if err := mizugos.Configmgr().Unmarshal(nameAuth, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameAuth, err)
	} // if

	this.listenID = mizugos.Netmgr().AddListenTCP(this.config.IP, this.config.Port, this.bind, this.unbind, this.listenWrong)
	mizugos.Info(defines.LogSystem, nameAuth).Caller(0).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Auth) Finalize() {
	mizugos.Netmgr().DelListen(this.listenID)
	mizugos.Info(defines.LogSystem, nameAuth).Caller(0).Message("entry finalize").End()
}

// bind 綁定處理
func (this *Auth) bind(session nets.Sessioner) *nets.Bundle {
	mizugos.Info(defines.LogSystem, nameAuth).Caller(0).Message("bind").End()
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

	if err := entity.SetProcess(procs.NewJson()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewAuth()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, nameAuth)
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	mizugos.Error(defines.LogSystem, nameAuth).Caller(0).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Auth) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// listenWrong 監聽錯誤處理
func (this *Auth) listenWrong(err error) {
	mizugos.Error(defines.LogSystem, nameAuth).Caller(1).EndError(err)
}

// bindWrong 綁定錯誤處理
func (this *Auth) bindWrong(err error) {
	mizugos.Warn(defines.LogSystem, nameAuth).Caller(1).EndError(err)
}
