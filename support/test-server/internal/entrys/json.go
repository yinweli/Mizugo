package entrys //nolint:dupl

import (
	"fmt"
	"sync/atomic"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
	"github.com/yinweli/Mizugo/support/test-server/internal/features"
	"github.com/yinweli/Mizugo/support/test-server/internal/modules"
)

// NewJson 建立Json入口
func NewJson() *Json {
	return &Json{
		name: "json",
	}
}

// Json Json入口
type Json struct {
	name     string        // 系統名稱
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
	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	this.listenID = mizugos.Netmgr().AddListenTCP(this.config.IP, this.config.Port, this.bind, this.unbind, this.listenWrong)
	features.LogSystem.Get().Info(this.name).Message("entry start").KV("config", this.config).Caller(0).EndFlush()
	return nil
}

// Finalize 結束處理
func (this *Json) Finalize() {
	mizugos.Netmgr().DelListen(this.listenID)
	features.LogSystem.Get().Info(this.name).Message("entry finalize").Caller(0).EndFlush()
}

// bind 綁定處理
func (this *Json) bind(session nets.Sessioner) *nets.Bundle {
	entity := mizugos.Entitymgr().Add()

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

	if err := entity.SetProcess(procs.NewJson().Base64(true).DesCBC(true, this.config.Key, this.config.Key)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.SetSession(session); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.AddModule(modules.NewJson(this.count.Add)); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	if err := entity.Initialize(this.bindWrong); err != nil {
		wrong = fmt.Errorf("bind: %w", err)
		goto Error
	} // if

	mizugos.Labelmgr().Add(entity, this.name)
	session.SetOwner(entity)
	features.LogSystem.Get().Info(this.name).Message("bind").Caller(0).EndFlush()
	return entity.Bundle()

Error:
	if entity != nil {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	} // if

	session.Stop()
	features.LogSystem.Get().Error(this.name).Caller(0).Error(wrong).EndFlush()
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
	features.LogSystem.Get().Error(this.name).Caller(1).Error(err).EndFlush()
}

// bindWrong 綁定錯誤處理
func (this *Json) bindWrong(err error) {
	features.LogSystem.Get().Warn(this.name).Caller(1).Error(err).EndFlush()
}
