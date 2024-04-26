package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
)

// MetricsInitialize 初始化度量
func MetricsInitialize() error {
	config := &MetricsConfig{}

	if err := mizugos.Config.Unmarshal("metrics", config); err != nil {
		return fmt.Errorf("metrics initialize: %w", err)
	} // if

	if err := mizugos.Metrics.Initialize(config.Port); err != nil {
		return fmt.Errorf("metrics initialize: %w", err)
	} // if

	MeterLogin = mizugos.Metrics.NewRuntime("login")
	MeterUpdate = mizugos.Metrics.NewRuntime("update")
	MeterJson = mizugos.Metrics.NewRuntime("json")
	MeterProto = mizugos.Metrics.NewRuntime("proto")
	MeterRaven = mizugos.Metrics.NewRuntime("raven")
	LogSystem.Get().Info("metrics").Message("initialize").EndFlush()
	return nil
}

// MetricsConfig 配置資料
type MetricsConfig struct {
	Port int `yaml:"port"` // 埠號
}

var MeterLogin *metrics.Runtime  // Login訊息度量物件
var MeterUpdate *metrics.Runtime // Update訊息度量物件
var MeterJson *metrics.Runtime   // Json訊息度量物件
var MeterProto *metrics.Runtime  // Proto訊息度量物件
var MeterRaven *metrics.Runtime  // Raven訊息度量物件
