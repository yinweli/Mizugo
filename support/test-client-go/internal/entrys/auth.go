package entrys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/modules"
)

// InitializeAuth 初始化Auth
func InitializeAuth() (err error) {
	if err = features.Config.Unmarshal(auth.name, &auth.config); err != nil {
		return fmt.Errorf("%v initialize: %w", auth.name, err)
	} // if

	if auth.config.Enable {
		features.Net.AddConnectTCP(auth.config.IP, auth.config.Port, auth.config.Timeout, auth.bind, auth.unbind, auth.connectWrong)
	} // if

	features.LogSystem.Get().Info(auth.name).Message("initialize").EndFlush()
	return nil
}

// FinalizeAuth 結束Auth
func FinalizeAuth() {
	features.LogSystem.Get().Info(auth.name).Message("finalize").EndFlush()
}

// Auth Auth入口
type Auth struct {
	name   string     // 系統名稱
	config AuthConfig // 配置資料
}

// AuthConfig 配置資料
type AuthConfig struct {
	Enable  bool          `yaml:"enable"`  // 啟用旗標
	IP      string        `yaml:"ip"`      // 位址
	Port    string        `yaml:"port"`    // 埠號
	Timeout time.Duration `yaml:"timeout"` // 超期時間
	Delay   time.Duration `yaml:"delay"`   // 延遲時間
	Account string        `yaml:"account"` // 帳號
	Update  int           `yaml:"update"`  // 更新次數
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

	if err := entity.AddModule(modules.NewAuth(this.config.Delay, this.config.Account, this.config.Update)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	features.Label.Add(entity, this.name)
	session.SetOwner(entity)
	features.MeterConnect.Add(1)
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
		features.MeterConnect.Add(-1)
	} // if
}

// connectWrong 連接錯誤處理
func (this *Auth) connectWrong(err error) {
	features.LogSystem.Get().Error(this.name).Caller(1).Error(err).EndFlush()
}

// bindWrong 綁定錯誤處理
func (this *Auth) bindWrong(err error) {
	features.LogSystem.Get().Warn(this.name).Caller(1).Error(err).EndFlush()
}

var auth = &Auth{name: "auth"} // Auth入口
