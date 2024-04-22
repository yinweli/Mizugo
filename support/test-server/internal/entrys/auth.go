package entrys

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// AuthInitialize 初始化Auth入口
func AuthInitialize() (err error) {
	config := &AuthConfig{}

	if err = mizugos.Config.Unmarshal("auth", config); err != nil {
		return fmt.Errorf("auth initialize: %w", err)
	} // if

	auth.listenID = mizugos.Network.AddListenTCP(config.IP, config.Port, auth.bind, auth.unbind, auth.wrong)
	features.LogSystem.Get().Info("auth").Message("initialize").EndFlush()
	return nil
}

// Auth Auth入口
type Auth struct {
	listenID nets.ListenID // 接聽編號
}

// AuthConfig 配置資料
type AuthConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
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

	if err = entity.AddModule(modules.NewAuth()); err != nil {
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
	} // if
}

// wrong 錯誤處理
func (this *Auth) wrong(err error) {
	features.LogSystem.Get().Error("auth").Caller(1).Error(err).EndFlush()
}

var auth = &Auth{} // Auth入口
