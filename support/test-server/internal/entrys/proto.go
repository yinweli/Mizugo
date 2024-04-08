package entrys //nolint:dupl

import (
	"fmt"
	"sync/atomic"

	"github.com/yinweli/Mizugo/mizugo/entitys"
	"github.com/yinweli/Mizugo/mizugo/nets"
	"github.com/yinweli/Mizugo/mizugo/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// InitializeProto 初始化Proto
func InitializeProto() error {
	if err := features.Config.Unmarshal(proto.name, &proto.config); err != nil {
		return fmt.Errorf("%v initialize: %w", proto.name, err)
	} // if

	proto.listenID = features.Net.AddListenTCP(proto.config.IP, proto.config.Port, proto.bind, proto.unbind, proto.listenWrong)
	features.LogSystem.Get().Info(proto.name).Message("initialize").EndFlush()
	return nil
}

// FinalizeProto 結束Proto
func FinalizeProto() {
	features.Net.DelListen(proto.listenID)
	features.LogSystem.Get().Info(proto.name).Message("finalize").EndFlush()
}

// Proto Proto資料
type Proto struct {
	name     string        // 系統名稱
	config   ProtoConfig   // 配置資料
	count    atomic.Int64  // 計數器
	listenID nets.ListenID // 接聽編號
}

// ProtoConfig 配置資料
type ProtoConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// bind 綁定處理
func (this *Proto) bind(session nets.Sessioner) *nets.Bundle {
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

	if err := entity.SetProcess(procs.NewProto().Base64(true).DesCBC(true, this.config.Key, this.config.Key)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewProto(this.count.Add)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
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
func (this *Proto) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		features.Entity.Del(entity.EntityID())
		features.Label.Erase(entity)
	} // if
}

// listenWrong 監聽錯誤處理
func (this *Proto) listenWrong(err error) {
	features.LogSystem.Get().Error(this.name).Caller(1).Error(err).EndFlush()
}

// bindWrong 綁定錯誤處理
func (this *Proto) bindWrong(err error) {
	features.LogSystem.Get().Warn(this.name).Caller(1).Error(err).EndFlush()
}

var proto = &Proto{name: "proto"} // Proto入口
