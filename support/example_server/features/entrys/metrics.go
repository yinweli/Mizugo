package entrys

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/support/example_server/features/defines"
)

// NewMetrics 建立統計入口資料
func NewMetrics() *Metrics {
	return &Metrics{
		name: defines.EntryMetrics,
	}
}

// Metrics 統計入口資料
type Metrics struct {
	name   string        // 入口名稱
	config MetricsConfig // 設定資料
}

// MetricsConfig 設定資料
type MetricsConfig struct {
	Port     int    // 埠號
	Username string // 帳號
	Password string // 密碼
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	mizugos.Info(this.name).Message("entry initialize").End()

	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	metrics.Initialize(this.config.Port, &metrics.Auth{
		Username: this.config.Username,
		Password: this.config.Password,
	})
	mizugos.Info(this.name).Message("entry start").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Info(this.name).Message("entry finalize").End()
	metrics.Finalize()
}
