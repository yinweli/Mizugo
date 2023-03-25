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

const nameJson = "json" // 入口名稱

// NewJson 建立Json入口
func NewJson() *Json {
	return &Json{}
}

// Json Json入口
type Json struct {
	config   JsonConfig    // 配置資料
	count    atomic.Int64  // 計數器
	listenID nets.ListenID // 接聽編號
}

// JsonConfig 配置資料
type JsonConfig struct {
	IP   string `yaml:"ip"`   // 位址
	Port string `yaml:"port"` // 埠號
	Key  string `yaml:"key"`  // 密鑰
}

// Initialize 初始化處理
func (this *Json) Initialize() error {
	features.System.Info(nameJson).Caller(0).Message("entry initialize").End()

	if err := mizugos.Configmgr().Unmarshal(nameJson, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameJson, err)
	} // if

	this.listenID = mizugos.Netmgr().AddListenTCP(this.config.IP, this.config.Port, this.bind, this.unbind, this.listenWrong)
	features.System.Info(nameJson).Caller(0).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Json) Finalize() {
	mizugos.Netmgr().DelListen(this.listenID)
	features.System.Info(nameJson).Caller(0).Message("entry finalize").End()
}

// bind 綁定處理
func (this *Json) bind(session nets.Sessioner) *nets.Bundle {
	features.System.Info(nameJson).Caller(0).Message("bind").End()
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

	if err := entity.SetProcess(procs.NewJson().Base64(true).DesCBC(true, this.config.Key, this.config.Key)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewJson(this.incr)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, nameJson)
	session.SetOwner(entity)
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	features.System.Error(nameJson).Caller(0).EndError(wrong)
	return nil
}

// unbind 解綁處理
func (this *Json) unbind(session nets.Sessioner) {
	if entity, ok := session.GetOwner().(*entitys.Entity); ok {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if
}

// listenWrong 監聽錯誤處理
func (this *Json) listenWrong(err error) {
	features.System.Error(nameJson).Caller(1).EndError(err)
}

// bindWrong 綁定錯誤處理
func (this *Json) bindWrong(err error) {
	features.System.Warn(nameJson).Caller(1).EndError(err)
}

// incr 增加計數
func (this *Json) incr() int64 {
	return this.count.Add(1)
}
