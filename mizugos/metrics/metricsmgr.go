package metrics

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/pools"
)

// NewMetricsmgr 建立度量管理器
func NewMetricsmgr() *Metricsmgr {
	return &Metricsmgr{}
}

// Metricsmgr 度量管理器, 其中包括兩部分: 效能數據(來自pprof), 自訂統計或統計數據(來自expvar)
//
// 如果想查看效能數據, 可以參考以下網址
//   - https://blog.csdn.net/skh2015java/article/details/102748222
//   - http://www.zyiz.net/tech/detail-112761.html
//   - https://www.iargs.cn/?p=62
//   - https://www.readfog.com/a/1635446409103773696
//
// 如果想建立自訂統計, 執行 NewInt, NewFloat, NewString, NewMap, 就可以獲得自訂記錄器;
// 如果想建立執行統計, 執行 NewRuntime 就可以獲得執行記錄器
//
// 如果想查看自訂統計或統計數據, 可以從 https://github.com/divan/expvarmon 安裝expvarmon工具;
// 工具參數說明如下:
//   - ports="http://網址:埠號"
//   - vars="監控名稱:記錄名稱,..."
//   - i 間隔時間
//
// 如果想用expvarmon工具查看執行統計數據, 假設執行記錄器的名稱為'echo', 則改變var參數為
// vars="time:echo.time,time(max):echo.time(max),time(avg):echo.time(avg),count:echo.count,count(1m):echo.count(1m),count(5m):echo.count(5m),count(10m):echo.count(10m),count(60m):echo.count(60m)"
//
// expvarmon範例如下:
//   - expvarmon -ports="http://localhost:8080" -i 1s
//   - expvarmon -ports="http://localhost:8080" -vars="count:count,total:total,money:value" -i 1s
type Metricsmgr struct {
	ctx    ctxs.Ctx     // ctx物件
	server *http.Server // http伺服器物件
}

// Initialize 初始化處理
func (this *Metricsmgr) Initialize(port int) error {
	if this.server != nil {
		return fmt.Errorf("metricsmgr initialize: already initialize")
	} // if

	handler := http.NewServeMux()
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	handler.Handle("/debug/vars", expvar.Handler())

	this.ctx = ctxs.Get().WithCancel()
	this.server = &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		ReadHeaderTimeout: timeout,
		Handler:           handler,
	}

	pools.DefaultPool.Submit(func() {
		_ = this.server.ListenAndServe()
	})
	return nil
}

// Finalize 結束處理
func (this *Metricsmgr) Finalize() {
	if this.server != nil {
		_ = this.server.Close()
		this.server = nil
	} // if

	this.ctx.Cancel()
}

// NewInt 建立整數統計
func (this *Metricsmgr) NewInt(name string) *Int {
	return expvar.NewInt(name)
}

// NewFloat 建立浮點數統計
func (this *Metricsmgr) NewFloat(name string) *Float {
	return expvar.NewFloat(name)
}

// NewString 建立字串統計
func (this *Metricsmgr) NewString(name string) *String {
	return expvar.NewString(name)
}

// NewMap 建立映射統計
func (this *Metricsmgr) NewMap(name string) *Map {
	return expvar.NewMap(name)
}

// NewRuntime 建立執行統計
func (this *Metricsmgr) NewRuntime(name string) *Runtime {
	v := &Runtime{}
	v.start(this.ctx.Ctx())
	expvar.Publish(name, v)
	return v
}
