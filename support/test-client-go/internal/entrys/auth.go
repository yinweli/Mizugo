package entrys

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/modules"
)

// AuthInitialize 初始化Auth入口
func AuthInitialize() (err error) {
	config := &AuthConfig{}

	if err = mizugos.Config.Unmarshal("auth", config); err != nil {
		return fmt.Errorf("auth initialize: %w", err)
	} // if

	auth.delay = config.Delay
	auth.account = config.Account
	auth.update = config.Update

	if config.Enable {
		mizugos.Network.AddConnectTCP(config.IP, config.Port, config.Timeout, auth.bind, auth.unbind, auth.wrong)
	} // if

	features.LogSystem.Get().Info("auth").Message("initialize").EndFlush()
	return nil
}

// Auth Auth入口
type Auth struct {
	delay   time.Duration // 延遲時間
	account string        // 帳號
	update  int           // 更新次數
}

// AuthConfig 配置資料
type AuthConfig struct {
	Enable  bool          `yaml:"enable"`  // 啟用旗標
	IP      string        `yaml:"ip"`      // 位址
	Port    string        `yaml:"port"`    // 埠號
	Timeout time.Duration `yaml:"timeout"` // 超時時間
	Delay   time.Duration `yaml:"delay"`   // 延遲時間
	Account string        `yaml:"account"` // 帳號
	Update  int           `yaml:"update"`  // 更新次數
}

// bind 綁定處理
func (this *Auth) bind(session nets.Sessioner) bool {
	err := error(nil)
	entity := mizugos.Entity.Add()
	process := procs.NewJson()

	if entity == nil {
		err = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err = entity.SetProcess(process); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.SetSession(session); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.AddModule(modules.NewAuth(this.delay, this.account, this.update)); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.Initialize(this.wrong); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	session.SetCodec(process)
	session.SetPublish(entity.PublishOnce)
	session.SetWrong(this.wrong)
	session.SetOwner(entity)
	features.MeterConnect.Add(1)
	features.LogSystem.Get().Info("auth").Message("bind").Caller(0).EndFlush()
	return true

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("auth").Caller(0).Error(err).EndFlush()
	return false
}

// unbind 解綁處理
func (this *Auth) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
		features.MeterConnect.Add(-1)
	} // if
}

// wrong 錯誤處理
func (this *Auth) wrong(err error) {
	features.LogSystem.Get().Error("auth").Caller(1).Error(err).EndFlush()
}

var auth = &Auth{} // Auth入口
