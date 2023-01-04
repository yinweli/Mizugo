package commons

import (
	"fmt"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/support/example_server/internal/defines"
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
	Port     int    // 埠號
	Username string // 帳號
	Password string // 密碼
}

// Initialize 初始化處理
func (this *Metrics) Initialize() error {
	if err := mizugos.Configmgr().ReadFile(this.name, defines.ConfigType); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	if err := mizugos.Configmgr().Unmarshal(this.name, &this.config); err != nil {
		return fmt.Errorf("%v initialize: %w", this.name, err)
	} // if

	mizugos.Metricsmgr().Initialize(this.config.Port, &metrics.Auth{
		Username: this.config.Username,
		Password: this.config.Password,
	})
	mizugos.Info(this.name).Message("initialize").KV("config", this.config).End()
	Echo = mizugos.Metricsmgr().NewRuntime("echo")
	return nil
}

// Finalize 結束處理
func (this *Metrics) Finalize() {
	mizugos.Metricsmgr().Finalize()
}

// Echo 回音統計物件. 使用expvarmon監控時, 可使用以下參數
// -ports="http://帳號:密碼@網址:埠號"
// -vars="time:echo.time,max:echo.max,mean:echo.mean,count:echo.count,count(1m):echo.count(1m),count(5m):echo.count(5m),count(10m):echo.count(10m),count(60m):echo.count(60m)"
var Echo *metrics.Runtime
