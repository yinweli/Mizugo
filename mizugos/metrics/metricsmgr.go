package metrics

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync/atomic"
)

// 指標管理器, 其中包括兩部分
//   效能指標(來自pprof)
//   自訂指標或統計數據(來自expvar)
// 如果想查看效能指標, 可以參考以下網址
//   https://blog.csdn.net/skh2015java/article/details/102748222
//   http://www.zyiz.net/tech/detail-112761.html
//   https://www.iargs.cn/?p=62
//   https://www.readfog.com/a/1635446409103773696
// 如果想查看自訂指標或統計數據, 可以通過以下工具
//   https://github.com/divan/expvarmon
// 此工具同時也可以查看記憶體使用情況, 可使用以下參數
//   -ports="http://網址:埠號"
//   -i 間隔時間
//   範例: expvarmon -ports="http://localhost:8080" -i 1s
//   範例: expvarmon -ports="http://localhost:8080" -vars="...自訂指標..." -i 1s
// 指標管理器同時還提供執行統計工具, 只要建立 Metricsmgr.NewRuntime(統計名稱) 就可以記錄特定區段的執行數據
// 如果要用expvarmon查看執行數據, 可以添加以下參數
//   假設執行數據的名稱為 'echo'
//   -vars="time:echo.time,max:echo.max,mean:echo.mean,count:echo.count,count(1m):echo.count(1m),count(5m):echo.count(5m),count(10m):echo.count(10m),count(60m):echo.count(60m)"

// NewMetricsmgr 建立指標管理器
func NewMetricsmgr() *Metricsmgr {
	return &Metricsmgr{}
}

// Metricsmgr 指標管理器
type Metricsmgr struct {
	stop   atomic.Bool  // 結束旗標
	server *http.Server // http伺服器物件
}

// Initialize 初始化處理
func (this *Metricsmgr) Initialize(port int) {
	handler := http.NewServeMux()
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	handler.Handle("/debug/vars", expvar.Handler())

	this.server = &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		ReadHeaderTimeout: serverTimeout,
		Handler:           handler,
	}

	go func() {
		_ = this.server.ListenAndServe()
	}()
}

// Finalize 結束處理
func (this *Metricsmgr) Finalize() {
	this.stop.Store(true)
	_ = this.server.Close()
}

// NewInt 建立整數統計
func (this *Metricsmgr) NewInt(name string) *expvar.Int {
	return expvar.NewInt(name)
}

// NewFloat 建立浮點數統計
func (this *Metricsmgr) NewFloat(name string) *expvar.Float {
	return expvar.NewFloat(name)
}

// NewString 建立字串統計
func (this *Metricsmgr) NewString(name string) *expvar.String {
	return expvar.NewString(name)
}

// NewMap 建立映射統計
func (this *Metricsmgr) NewMap(name string) *expvar.Map {
	return expvar.NewMap(name)
}

// NewRuntime 建立執行統計
func (this *Metricsmgr) NewRuntime(name string) *Runtime {
	v := &Runtime{
		finish: func() bool {
			return this.stop.Load()
		},
	}
	v.start()
	expvar.Publish(name, v)
	return v
}
