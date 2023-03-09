package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/support/test-server/internal/defines"
)

// NewMetrics 建立統計資料
func NewMetrics() *Metrics {
	return &Metrics{
		name: "metrics",
	}
}

// Metrics 統計資料
type Metrics struct {
	name   string        // 特性名稱
	config MetricsConfig // 配置資料
}

// MetricsConfig 配置資料
type MetricsConfig struct {
	Port int `yaml:"port"` // 埠號
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Metricsmgr().Initialize(this.config.Port); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	Login = mizugos.Metricsmgr().NewRuntime("login")
	Update = mizugos.Metricsmgr().NewRuntime("update")
	Json = mizugos.Metricsmgr().NewRuntime("json")
	Proto = mizugos.Metricsmgr().NewRuntime("proto")
	mizugos.Info(defines.LogSystem, this.name).Caller(0).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}

var Login *metrics.Runtime  // Login訊息統計物件
var Update *metrics.Runtime // Update訊息統計物件
var Json *metrics.Runtime   // Json訊息統計物件
var Proto *metrics.Runtime  // Proto訊息統計物件
