package entrys //nolint:dupl

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// RavenInitialize 初始化Raven入口
func RavenInitialize() (err error) {
	config := &RavenConfig{}

	if err = mizugos.Config.Unmarshal("raven", config); err != nil {
		return fmt.Errorf("raven initialize: %w", err)
	} // if

	raven.listenID = mizugos.Network.AddListenTCP(config.IP, config.Port, raven.bind, raven.unbind, raven.wrong)
	raven.key = config.Key
	features.LogSystem.Get().Info("raven").Message("initialize").EndFlush()
	return nil
}

// Raven Raven入口
type Raven struct {
	listenID nets.ListenID // 接聽編號
	key      string        // 密鑰
}

// RavenConfig 配置資料
type RavenConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// bind 綁定處理
func (this *Raven) bind(session nets.Sessioner) bool {
	err := error(nil)
	entity := mizugos.Entity.Add()
	process := procs.NewRaven()

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

	if err = entity.AddModule(modules.NewRaven()); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.Initialize(this.wrong); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	session.SetCodec(process, cryptos.NewDesCBC(cryptos.PaddingPKCS7, this.key, this.key), cryptos.NewBase64())
	session.SetPublish(entity.PublishOnce)
	session.SetWrong(this.wrong)
	session.SetOwner(entity)
	features.LogSystem.Get().Info("raven").Message("bind").Caller(0).EndFlush()
	return true

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("raven").Caller(0).Error(err).EndFlush()
	return false
}

// unbind 解綁處理
func (this *Raven) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if
}

// wrong 錯誤處理
func (this *Raven) wrong(err error) {
	features.LogSystem.Get().Error("raven").Caller(1).Error(err).EndFlush()
}

var raven = &Raven{} // Raven入口
