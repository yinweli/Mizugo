package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
)

// NewMetrics 建立統計資料
func NewMetrics() *Metrics {
	return &Metrics{
		name: "metrics",
	}
}

// Metrics 統計資料
type Metrics struct {
	name   string        // 系統名稱
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

	MeterLogin = mizugos.Metricsmgr().NewRuntime("login")
	MeterUpdate = mizugos.Metricsmgr().NewRuntime("update")
	MeterJson = mizugos.Metricsmgr().NewRuntime("json")
	MeterProto = mizugos.Metricsmgr().NewRuntime("proto")
	LogSystem.Get().Info(this.name).Message("initialize").KV("config", this.config).Caller(0).EndFlush()
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}

var MeterLogin *metrics.Runtime  // Login訊息統計物件
var MeterUpdate *metrics.Runtime // Update訊息統計物件
var MeterJson *metrics.Runtime   // Json訊息統計物件
var MeterProto *metrics.Runtime  // Proto訊息統計物件
