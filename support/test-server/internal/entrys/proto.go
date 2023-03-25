package entrys //nolint:dupl

import (
	"fmt"
	"sync/atomic"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

const nameProto = "proto" // 入口名稱

// NewProto 建立Proto入口
func NewProto() *Proto {
	return &Proto{}
}

// Proto Proto入口
type Proto struct {
	config   ProtoConfig   // 配置資料
	listenID nets.ListenID // 接聽編號
	count    atomic.Int64  // 計數器
}

// ProtoConfig 配置資料
type ProtoConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// Initialize 初始化處理
func (this *Proto) Initialize() error {
	features.System.Info(nameProto).Caller(0).Message("entry initialize").End()

	if err := mizugos.Configmgr().Unmarshal(nameProto, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameProto, err)
	} // if

	this.listenID = mizugos.Netmgr().AddListenTCP(this.config.IP, this.config.Port, this.bind, this.unbind, this.listenWrong)
	features.System.Info(nameProto).Caller(0).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Proto) Finalize() {
	mizugos.Netmgr().DelListen(this.listenID)
	features.System.Info(nameProto).Caller(0).Message("entry finalize").End()
}

// bind 綁定處理
func (this *Proto) bind(session nets.Sessioner) *nets.Bundle {
	features.System.Info(nameProto).Caller(0).Message("bind").End()
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

	if err := entity.SetProcess(procs.NewProto().Base64(true).DesCBC(true, this.config.Key, this.config.Key)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewProto(this.incr)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, nameProto)
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	features.System.Error(nameProto).Caller(0).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Proto) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// listenWrong 監聽錯誤處理
func (this *Proto) listenWrong(err error) {
	features.System.Error(nameProto).Caller(1).EndError(err)
}

// bindWrong 綁定錯誤處理
func (this *Proto) bindWrong(err error) {
	features.System.Warn(nameProto).Caller(1).EndError(err)
}

// incr 增加計數
func (this *Proto) incr() int64 {
	return this.count.Add(1)
}
