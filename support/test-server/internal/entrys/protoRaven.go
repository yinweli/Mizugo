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

// ProtoRavenInitialize 初始化ProtoRaven入口
func ProtoRavenInitialize() (err error) {
	config := &ProtoRavenConfig{}

	if err = mizugos.Config.Unmarshal("protoRaven", config); err != nil {
		return fmt.Errorf("protoRaven initialize: %w", err)
	} // if

	protoRaven.listenID = mizugos.Network.AddListenTCP(config.IP, config.Port, protoRaven.bind, protoRaven.unbind, protoRaven.wrong)
	protoRaven.key = config.Key
	features.LogSystem.Get().Info("protoRaven").Message("initialize").EndFlush()
	return nil
}

// ProtoRaven ProtoRaven入口
type ProtoRaven struct {
	listenID nets.ListenID // 接聽編號
	key      string        // 密鑰
}

// ProtoRavenConfig 配置資料
type ProtoRavenConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// bind 綁定處理
func (this *ProtoRaven) bind(session nets.Sessioner) bool {
	err := error(nil)
	entity := mizugos.Entity.Add()
	process := procs.NewProto()

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

	if err = entity.AddModule(modules.NewProtoRaven()); err != nil {
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
	features.LogSystem.Get().Info("protoRaven").Message("bind").Caller(0).EndFlush()
	return true

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("protoRaven").Caller(0).Error(err).EndFlush()
	return false
}

// unbind 解綁處理
func (this *ProtoRaven) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if
}

// wrong 錯誤處理
func (this *ProtoRaven) wrong(err error) {
	features.LogSystem.Get().Error("protoRaven").Caller(1).Error(err).EndFlush()
}

var protoRaven = &ProtoRaven{} // ProtoRaven入口
