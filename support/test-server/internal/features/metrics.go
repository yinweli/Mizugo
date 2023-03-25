package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
)

const nameMetrics = "metrics" // 特性名稱

// NewMetrics 建立統計資料
func NewMetrics() *Metrics {
	return &Metrics{}
}

// Metrics 統計資料
type Metrics struct {
	config MetricsConfig // 配置資料
}

// MetricsConfig 配置資料
type MetricsConfig struct {
	Port int `yaml:"port"` // 埠號
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(nameMetrics, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameMetrics, err)
	} // if

	if err := mizugos.Metricsmgr().Initialize(this.config.Port); err != nil {
		return fmt.Errorf("%v initialize: %w", nameMetrics, err)
	} // if

	Login = mizugos.Metricsmgr().NewRuntime("login")
	Update = mizugos.Metricsmgr().NewRuntime("update")
	Json = mizugos.Metricsmgr().NewRuntime("json")
	Proto = mizugos.Metricsmgr().NewRuntime("proto")
	System.Info(nameMetrics).Caller(0).Message("initialize").KV("config", this.config).End()
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
