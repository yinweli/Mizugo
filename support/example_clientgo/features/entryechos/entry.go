package entryechos

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/cores/msgs"
	"github.com/yinweli/Mizugo/cores/nets"
	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/features/defines"
)

// NewEntry 建立入口資料
func NewEntry() *Entry {
	return &Entry{
		name: "echoc",
	}
}

// Entry 入口資料
type Entry struct {
	name   string // 入口名稱
	config Config // 設定資料
}

// Config 設定資料
type Config struct {
	IP      string        // 位址
	Port    string        // 埠號
	Timeout time.Duration // 逾期時間(秒)
}

// Initialize 初始化處理
func (this *Entry) Initialize() error {
	mizugos.Info(this.name).
		Message("entry initialize").
		End()

	if err := mizugos.Configmgr().ReadFile(this.name, "yaml"); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Netmgr().AddConnect(nets.NewTCPConnect(this.config.IP, this.config.Port, this.config.Timeout), this)
	mizugos.Info(this.name).
		Message("entry start").
		KV("ip", this.config.IP).
		KV("port", this.config.Port).
		End()
	return nil
}

// Finalize 結束處理
func (this *Entry) Finalize() {
	mizugos.Info(this.name).
		Message("entry stop").
		End()
}

// Bind 綁定處理
func (this *Entry) Bind(session nets.Sessioner) (content nets.Content, err error) {
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(msgs.NewStringProc()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	// TODO: add module

	if err := entity.Initialize(func() {
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	}); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEcho)
	content.Unbind = entity.Finalize
	content.Encode = nil
	content.Decode = nil
	content.Receive = nil
	return content, nil
}

// Error 錯誤處理
func (this *Entry) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
