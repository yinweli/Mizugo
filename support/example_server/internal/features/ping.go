package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
	"github.com/yinweli/Mizugo/support/example_server/internal/modules"
)

// NewPing 建立ping入口資料
func NewPing() *Ping {
	return &Ping{
		name: defines.EntryPing,
	}
}

// Ping ping入口資料
type Ping struct {
	name     string        // 入口名稱
	config   PingConfig    // 設定資料
	listenID nets.ListenID // 接聽編號
}

// PingConfig 設定資料
type PingConfig struct {
	IP      string // 位址
	Port    string // 埠號
	InitKey string // 初始密鑰
}

// Initialize 初始化處理
func (this *Ping) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	this.listenID = mizugos.Netmgr().AddListen(nets.NewTCPListen(this.config.IP, this.config.Port), this)
	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Ping) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	mizugos.Netmgr().DelListen(this.listenID)
}

// Bind 綁定處理
func (this *Ping) Bind(session nets.Sessioner) (content nets.Content, err error) {
	mizugos.Info(this.name).Message("session").KV("sessionID", session.SessionID()).End()
	entity := mizugos.Entitymgr().Add()

	if entity == nil {
		return content, fmt.Errorf("bind: entity nil")
	} // if

	if err := entity.SetSession(session); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.SetProcess(procs.NewProtoDes([]byte(this.config.InitKey))); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	moduleKey := modules.NewKey()

	if err := entity.AddModule(moduleKey); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	modulePing := modules.NewPing()

	if err := entity.AddModule(modulePing); err != nil {
		return content, fmt.Errorf("bind: %w", err)
	} // if

	if err := entity.Initialize(); err != nil {
		mizugos.Entitymgr().Del(entity.EntityID())
		return content, fmt.Errorf("bind: %w", err)
	} // if

	mizugos.Labelmgr().Add(entity, defines.LabelPing)
	content.Unbind = func() {
		entity.Finalize()
		mizugos.Netmgr().DelSession(session.SessionID())
		mizugos.Entitymgr().Del(entity.EntityID())
		mizugos.Labelmgr().Erase(entity)
	}
	content.Encode = entity.GetProcess().Encode
	content.Decode = entity.GetProcess().Decode
	content.Receive = entity.GetProcess().Process
	content.AfterSend = moduleKey.AfterSend
	return content, nil
}

// TODO: 可能要測試一下主動中斷客戶端是否會造成錯誤

// Error 錯誤處理
func (this *Ping) Error(err error) {
	_ = mizugos.Error(this.name).EndError(err)
}
