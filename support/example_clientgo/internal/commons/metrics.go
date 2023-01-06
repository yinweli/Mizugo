package commons

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/defines"
)

// NewMetrics 建立統計資料
func NewMetrics() *Metrics {
	return &Metrics{
		name: "metrics",
	}
}

// Metrics 統計資料
type Metrics struct {
	name   string        // 統計名稱
	config MetricsConfig // 設定資料
}

// MetricsConfig 設定資料
type MetricsConfig struct {
	Port int // 埠號
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Metricsmgr().Initialize(this.config.Port)
	mizugos.Info(this.name).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}
