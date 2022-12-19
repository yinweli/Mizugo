package echos

import (
	"fmt"
	"path/filepath"

	"github.com/yinweli/Mizugo/cores/nets"
	"github.com/yinweli/Mizugo/mizugos"
)

// NewServer 建立伺服器資料
func NewServer(configPath string) *Server {
	return &Server{
		serverName: "echo",
		configName: "echo",
		configFile: "echo.yaml",
		configPath: configPath,
	}
}

// Server 伺服器資料
type Server struct {
	serverName string // 伺服器名稱
	configName string // 設定名稱
	configFile string // 設定檔案名稱
	configPath string // 設定檔案路徑
	config     config // 設定資料
}

// Config 設定資料
type config struct {
	IP   string `yaml:"ip"`   // 接聽位址
	Port string `yaml:"port"` // 接聽埠號
}

// Initialize 初始化處理
func (this *Server) Initialize() error {
	mizugos.Info(this.serverName).
		Message("server initialize").
		End()

	if err := mizugos.Configmgr().ReadFile(filepath.Join(this.configPath, this.configFile)); err != nil {
		return fmt.Errorf("%v initialize: %w", this.serverName, err)
	} // if

	if err := mizugos.Configmgr().GetObject(this.configName, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.serverName, err)
	} // if

	mizugos.Netmgr().AddListen(nets.NewTCPListen(this.config.IP, this.config.Port), this)
	mizugos.Info(this.serverName).
		Message("server start").
		KV("ip", this.config.IP).
		KV("port", this.config.Port).
		End()
	return nil
}

// Finalize 結束處理
func (this *Server) Finalize() {

}

// Bind 綁定處理
func (this *Server) Bind(session nets.Sessioner) (reactor nets.Reactor, unbinder nets.Unbinder) {
	return nil, nil
}

// Error 錯誤處理
func (this *Server) Error(err error) {
	_ = mizugos.Error(this.serverName).EndError(err)
}
