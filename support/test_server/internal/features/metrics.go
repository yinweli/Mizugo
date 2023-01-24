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
	name   string        // 統計名稱
	config MetricsConfig // 配置資料
}

// MetricsConfig 配置資料
type MetricsConfig struct {
	Port int // 埠號
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Metricsmgr().Initialize(this.config.Port); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Info(this.name).Message("initialize").KV("config", this.config).End()
	Key = mizugos.Metricsmgr().NewRuntime("key")
	Json = mizugos.Metricsmgr().NewRuntime("json")
	Proto = mizugos.Metricsmgr().NewRuntime("proto")
	Stack = mizugos.Metricsmgr().NewRuntime("stack")
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}

var Key *metrics.Runtime   // 密鑰訊息統計物件
var Json *metrics.Runtime  // json訊息統計物件
var Proto *metrics.Runtime // proto訊息統計物件
var Stack *metrics.Runtime // stack訊息統計物件
