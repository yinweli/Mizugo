package features

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/defines"
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
	Port int // 埠號
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	if err := mizugos.Configmgr().Unmarshal(nameMetrics, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", nameMetrics, err)
	} // if

	if err := mizugos.Metricsmgr().Initialize(this.config.Port); err != nil {
		return fmt.Errorf("%v initialize: %w", nameMetrics, err)
	} // if

	Auth = mizugos.Metricsmgr().NewRuntime("auth")
	Json = mizugos.Metricsmgr().NewRuntime("json")
	Proto = mizugos.Metricsmgr().NewRuntime("proto")
	Connect = mizugos.Metricsmgr().NewInt("connect")
	mizugos.Info(defines.LogSystem, nameMetrics).Caller(0).Message("initialize").KV("config", this.config).End()
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}

var Auth *metrics.Runtime  // auth訊息統計物件
var Json *metrics.Runtime  // json訊息統計物件
var Proto *metrics.Runtime // proto訊息統計物件
var Connect *metrics.Int   // 連線統計物件
