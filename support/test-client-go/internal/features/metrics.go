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

	MeterAuth = mizugos.Metricsmgr().NewRuntime("auth")
	MeterJson = mizugos.Metricsmgr().NewRuntime("json")
	MeterProto = mizugos.Metricsmgr().NewRuntime("proto")
	MeterConnect = mizugos.Metricsmgr().NewInt("connect")
	LogSystem.Info(this.name).Caller(0).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}

var MeterAuth *metrics.Runtime  // auth訊息統計物件
var MeterJson *metrics.Runtime  // json訊息統計物件
var MeterProto *metrics.Runtime // proto訊息統計物件
var MeterConnect *metrics.Int   // 連線統計物件
