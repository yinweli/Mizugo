package entrys

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// InitializeAuth 初始化Auth
func InitializeAuth() (err error) {
	if err = features.Config.Unmarshal(auth.name, &auth.config); err != nil {
		return fmt.Errorf("%v initialize: %w", auth.name, err)
	} // if

	auth.listenID = features.Net.AddListenTCP(auth.config.IP, auth.config.Port, auth.bind, auth.unbind, auth.listenWrong)
	features.LogSystem.Get().Info(auth.name).Message("initialize").EndFlush()
	return nil
}

// FinalizeAuth 結束Auth
func FinalizeAuth() {
	features.Net.DelListen(auth.listenID)
	features.LogSystem.Get().Info(auth.name).Message("finalize").EndFlush()
}

// Auth Auth入口
type Auth struct {
	name     string        // 名稱
	config   AuthConfig    // 配置資料
	listenID nets.ListenID // 接聽編號
}

// AuthConfig 配置資料
type AuthConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
}

// bind 綁定處理
func (this *Auth) bind(session nets.Sessioner) *nets.Bundle {
	entity := features.Entity.Add()

	var wrong error

	if entity == nil {
		wrong = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err := entity.SetModulemap(entitys.NewModulemap()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetEventmap(entitys.NewEventmap(defines.EventCapacity)); err != nil {
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

	if err := entity.Initialize(auth.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	features.Label.Add(entity, this.name)
	session.SetOwner(entity)
	features.LogSystem.Get().Info(this.name).Message("bind").Caller(0).EndFlush()
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		features.Entity.Del(entity.EntityID())
		features.Label.Erase(entity)
	} // if

	session.Stop()
	features.LogSystem.Get().Error(this.name).Caller(0).Error(wrong).EndFlush()
	return nil
}

// unbind 解綁處理
func (this *Auth) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		features.Entity.Del(entity.EntityID())
		features.Label.Erase(entity)
	} // if
}

// listenWrong 監聽錯誤處理
func (this *Auth) listenWrong(err error) {
	features.LogSystem.Get().Error(this.name).Caller(1).Error(err).EndFlush()
}

// bindWrong 綁定錯誤處理
func (this *Auth) bindWrong(err error) {
	features.LogSystem.Get().Warn(this.name).Caller(1).Error(err).EndFlush()
}

var auth = &Auth{name: "auth"} // Auth入口
