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

	MeterConnect = mizugos.Metrics.NewInt("connect")
	MeterAuth = mizugos.Metrics.NewRuntime("login")
	MeterJson = mizugos.Metrics.NewRuntime("json")
	MeterProto = mizugos.Metrics.NewRuntime("proto")
	MeterProtoRaven = mizugos.Metrics.NewRuntime("protoRaven")
	LogSystem.Get().Info("metrics").Message("initialize").EndFlush()
	return nil
}

// MetricsConfig 配置資料
type MetricsConfig struct {
	Port int `yaml:"port"` // 埠號
}

var MeterConnect *metrics.Int        // 連線度量物件
var MeterAuth *metrics.Runtime       // auth訊息度量物件
var MeterJson *metrics.Runtime       // json訊息度量物件
var MeterProto *metrics.Runtime      // proto訊息度量物件
var MeterProtoRaven *metrics.Runtime // proto訊息度量物件
