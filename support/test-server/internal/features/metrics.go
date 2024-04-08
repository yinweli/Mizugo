package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugo/metrics"
)

// InitializeMetrics 初始化度量管理器
func InitializeMetrics() error {
	name := "metrics"
	config := &MetricsConfig{}
	Metrics = metrics.NewMetricsmgr()

	if err := Config.Unmarshal(name, config); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	if err := Metrics.Initialize(config.Port); err != nil {
		return fmt.Errorf("%v initialize: %w", name, err)
	} // if

	MeterLogin = Metrics.NewRuntime("login")
	MeterUpdate = Metrics.NewRuntime("update")
	MeterJson = Metrics.NewRuntime("json")
	MeterProto = Metrics.NewRuntime("proto")
	LogSystem.Get().Info(name).Message("initialize").EndFlush()
	return nil
}

// FinalizeMetrics 結束度量管理器
func FinalizeMetrics() {
	if Metrics != nil {
		Metrics.Finalize()
	} // if
}

// MetricsConfig 配置資料
type MetricsConfig struct {
	Port int `yaml:"port"` // 埠號
}

var Metrics *metrics.Metricsmgr  // 度量管理器
var MeterLogin *metrics.Runtime  // Login訊息度量物件
var MeterUpdate *metrics.Runtime // Update訊息度量物件
var MeterJson *metrics.Runtime   // Json訊息度量物件
var MeterProto *metrics.Runtime  // Proto訊息度量物件
