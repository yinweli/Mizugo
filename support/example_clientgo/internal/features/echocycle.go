package features

import (
	"fmt"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/modules"
)

// NewEchoCycle 建立循環回音資料
func NewEchoCycle() *EchoCycle {
	return &EchoCycle{
		name: defines.EntryEchoCycle,
	}
}

// EchoCycle 循環回音資料
type EchoCycle struct {
	name   string          // 入口名稱
	config EchoCycleConfig // 設定資料
}

// EchoCycleConfig 設定資料
type EchoCycleConfig struct {
	IP         string        // 位址
	Port       string        // 埠號
	Timeout    time.Duration // 逾期時間(秒)
	Message    string        // 回音字串
	Disconnect bool          // 斷線旗標
	Reconnect  bool          // 重連旗標
	WaitTime   time.Duration // 重連等待時間
}

// Initialize 初始化處理
func (this *EchoCycle) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	go this.connect()
	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *EchoCycle) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
}

// Bind 綁定處理
func (this *EchoCycle) Bind(session nets.Sessioner) (content nets.Content, err error) {
	mizugos.Info(this.name).Message("session").KV("sessionID", session.SessionID()).End()
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(procs.NewSimple()); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.AddModule(modules.NewEchoCycle(this.config.Message, this.config.Disconnect)); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.Initialize(); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelEchoCycle)
	content.Unbind = func() {
		entity.Finalize()
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)

		if this.config.Reconnect { // TODO: 重連改成在外面跑執行緒檢查是否有需要重連
			time.Sleep(this.config.WaitTime)
			go this.connect()
		} // if
	}
	content.Encode = entity.GetProcess().Encode
	content.Decode = entity.GetProcess().Decode
	content.Receive = entity.GetProcess().Process
	return content, nil
}

// Error 錯誤處理
func (this *EchoCycle) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}

// connect 進行連線
func (this *EchoCycle) connect() {
	mizugos.Netmgr().AddConnect(nets.NewTCPConnect(this.config.IP, this.config.Port, this.config.Timeout), this)
}
