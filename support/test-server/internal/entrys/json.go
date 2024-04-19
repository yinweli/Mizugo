package entrys //nolint:dupl

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// JsonInitialize 初始化Json入口
func JsonInitialize() (err error) {
	config := &JsonConfig{}

	if err = mizugos.Config.Unmarshal("json", config); err != nil {
		return fmt.Errorf("json initialize: %w", err)
	} // if

	json.listenID = mizugos.Network.AddListenTCP(config.IP, config.Port, json.bind, json.unbind, json.wrong)
	json.key = config.Key
	features.LogSystem.Get().Info("json").Message("initialize").EndFlush()
	return nil
}

// Json Json入口
type Json struct {
	listenID nets.ListenID // 接聽編號
	key      string        // 密鑰
}

// JsonConfig 配置資料
type JsonConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// bind 綁定處理
func (this *Json) bind(session nets.Sessioner) bool {
	err := error(nil)
	entity := mizugos.Entity.Add()
	process := procs.NewJson()
	desCBC := cryptos.NewDesCBC(cryptos.PaddingPKCS7, this.key, this.key)
	base64 := cryptos.NewBase64()

	session.SetPublish(entity.PublishOnce)
	session.SetWrong(this.wrong)
	session.SetCodec(process, desCBC, base64)

	if entity == nil {
		err = fmt.Errorf("bind: entity nil")
		goto Error
	} // if

	if err = entity.SetModulemap(entitys.NewModulemap()); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.SetEventmap(entitys.NewEventmap(defines.EventCapacity)); err != nil {
		err = fmt.Errorf("bind: %w", err)
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

	if err = entity.AddModule(modules.NewJson()); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err = entity.Initialize(this.wrong); err != nil {
		err = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	session.SetOwner(entity)
	features.LogSystem.Get().Info("json").Message("bind").Caller(0).EndFlush()
	return true

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("json").Caller(0).Error(err).EndFlush()
	return false
}

// unbind 解綁處理
func (this *Json) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if
}

// wrong 錯誤處理
func (this *Json) wrong(err error) {
	features.LogSystem.Get().Error("json").Caller(1).Error(err).EndFlush()
}

var json = &Json{} // Json入口
