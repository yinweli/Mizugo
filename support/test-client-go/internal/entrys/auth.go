package entrys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/modules"
)

const nameAuth = "auth" // 入口名稱

// NewAuth 建立Auth入口
func NewAuth() *Auth {
	return &Auth{}
}

// Auth Auth入口
type Auth struct {
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

// Initialize 初始化處理
func (this *Auth) Initialize() error {
	features.System.Info(nameAuth).Caller(0).Message("entry initialize").End()

	if err := mizugos.Configmgr().Unmarshal(nameAuth, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameAuth, err)
	} // if

	if this.config.Enable {
		mizugos.Netmgr().AddConnectTCP(this.config.IP, this.config.Port, this.config.Timeout, this.bind, this.unbind, this.connectWrong)
	} // if

	features.System.Info(nameAuth).Caller(0).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Auth) Finalize() {
	features.System.Info(nameAuth).Caller(0).Message("entry finalize").End()
}

// bind 綁定處理
func (this *Auth) bind(session nets.Sessioner) *nets.Bundle {
	features.System.Info(nameAuth).Caller(0).Message("bind").End()
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

	if err := entity.AddModule(modules.NewAuth(this.config.Delay, this.config.Account, this.config.Update)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, nameAuth)
	features.Connect.Add(1)
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	features.System.Error(nameAuth).Caller(0).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Auth) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
		features.Connect.Add(-1)
	} // if
}

// connectWrong 連接錯誤處理
func (this *Auth) connectWrong(err error) {
	features.System.Error(nameAuth).Caller(1).EndError(err)
}

// bindWrong 綁定錯誤處理
func (this *Auth) bindWrong(err error) {
	features.System.Warn(nameAuth).Caller(1).EndError(err)
}
