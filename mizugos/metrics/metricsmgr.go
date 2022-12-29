package metrics

import (
	"expvar"
	"fmt"
	"net/http"
	"sync/atomic"
)

// TODO: 可能要考慮放棄用帳號密碼保護監測伺服器了

// NewMetricsmgr 建立統計管理器
func NewMetricsmgr() *Metricsmgr {
	return &Metricsmgr{}
}

// Metricsmgr 統計管理器
type Metricsmgr struct {
	stop   atomic.Bool  // 結束旗標
	server *http.Server // 監控伺服器物件
}

// Auth 認證資料
type Auth struct {
	Username string // 帳號
	Password string // 密碼
}

// Initialize 初始化處理, 如果有提供auth資料, 則會為監控伺服器添加帳號密碼認證
func (this *Metricsmgr) Initialize(port int, auth *Auth) {
	handler := http.NewServeMux()

	if auth != nil {
		handler.Handle(pattern, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if username, password, ok := r.BasicAuth(); ok == false || username != auth.Username || password != auth.Password {
					w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				} // if

				next.ServeHTTP(w, r)
			})
		}(expvar.Handler()))
	} else {
		handler.Handle(pattern, expvar.Handler())
	} // if

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
