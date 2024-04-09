package entrys //nolint:dupl

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// ProtoInitialize 初始化Proto入口
func ProtoInitialize() (err error) {
	config := &ProtoConfig{}

	if err = mizugos.Config.Unmarshal("proto", config); err != nil {
		return fmt.Errorf("proto initialize: %w", err)
	} // if

	proto.listenID = mizugos.Network.AddListenTCP(config.IP, config.Port, proto.bind, proto.unbind, proto.listenWrong)
	proto.key = config.Key
	features.LogSystem.Get().Info("proto").Message("initialize").EndFlush()
	return nil
}

// Proto Proto入口
type Proto struct {
	listenID nets.ListenID // 接聽編號
	key      string        // 密鑰
}

// ProtoConfig 配置資料
type ProtoConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// bind 綁定處理
func (this *Proto) bind(session nets.Sessioner) *nets.Bundle {
	entity := mizugos.Entity.Add()

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

	if err := entity.SetProcess(procs.NewProto().Base64(true).DesCBC(true, this.key, this.key)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewProto()); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	session.SetOwner(entity)
	features.LogSystem.Get().Info("proto").Message("bind").Caller(0).EndFlush()
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if

	session.Stop()
	features.LogSystem.Get().Error("proto").Caller(0).Error(wrong).EndFlush()
	return nil
}

// unbind 解綁處理
func (this *Proto) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entity.Del(entity.EntityID())
	} // if
}

// listenWrong 監聽錯誤處理
func (this *Proto) listenWrong(err error) {
	features.LogSystem.Get().Error("proto").Caller(1).Error(err).EndFlush()
}

// bindWrong 綁定錯誤處理
func (this *Proto) bindWrong(err error) {
	features.LogSystem.Get().Warn("proto").Caller(1).Error(err).EndFlush()
}

var proto = &Proto{} // Proto入口
