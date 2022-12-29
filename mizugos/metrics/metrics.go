package metrics

import (
	"expvar"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 監控統計數據的命令列工具
// https://tonybai.com/2021/04/14/expvarmon-save-and-convert-to-xlsx/
// https://github.com/divan/expvarmon

const (
	interval1  = 60   // 間隔時間: 1分鐘
	interval5  = 300  // 間隔時間: 5分鐘
	interval10 = 600  // 間隔時間: 10分鐘
	interval60 = 3600 // 間隔時間: 60分鐘
)

// Auth 認證資料
type Auth struct {
	Username string // 帳號
	Password string // 密碼
}

// Initialize 初始化處理, 如果有提供auth資料, 則會為監控的http伺服器添加帳號密碼認證
func Initialize(port int, auth *Auth) {
	go func() {
		stop.Store(false)
		server := &http.Server{
			Addr:              fmt.Sprintf(":%v", port),
			ReadHeaderTimeout: time.Second * 5,
			Handler:           handle(auth),
		}
		_ = server.ListenAndServe()
	}()
}

// Finalize 結束處理
func Finalize() {
	stop.Store(true)
}

// NewInt 產生整數統計
func NewInt(name string) *expvar.Int {
	return expvar.NewInt(name)
}

// NewFloat 產生浮點數統計
func NewFloat(name string) *expvar.Float {
	return expvar.NewFloat(name)
}

// NewString 產生字串統計
func NewString(name string) *expvar.String {
	return expvar.NewString(name)
}

// NewMap 產生映射統計
func NewMap(name string) *expvar.Map {
	return expvar.NewMap(name)
}

// NewRuntime 產生執行統計
func NewRuntime(name string) *Runtime {
	v := &Runtime{}
	v.start()
	expvar.Publish(name, v)
	return v
}

// Runtime 執行統計
type Runtime struct {
	stat runtime      // 統計資料
	curr runtime      // 當前資料
	lock sync.RWMutex // 執行緒鎖
}

// runtime 執行資料
type runtime struct {
	min     time.Duration // 最小執行時間
	max     time.Duration // 最大執行時間
	total   time.Duration // 總執行時間
	count   int64         // 總執行次數
	count1  int64         // 每分鐘執行次數
	count5  int64         // 每5分鐘執行次數
	count10 int64         // 每10分鐘執行次數
	count60 int64         // 每60分鐘執行次數
}

// Add 新增統計
func (this *Runtime) Add(delta time.Duration) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.curr.min > delta || this.curr.min == 0 {
		this.curr.min = delta
	} // if

	if this.curr.max < delta {
		this.curr.max = delta
	} // if

	this.curr.total += delta
	this.curr.count++
	this.curr.count1++
	this.curr.count5++
	this.curr.count10++
	this.curr.count60++
}

// Rec 新增統計
func (this *Runtime) Rec() func() {
	start := time.Now()
	return func() {
		this.Add(time.Since(start))
	}
}

// String 取得統計字串
func (this *Runtime) String() string {
	this.lock.RLock()
	stat := this.stat
	this.lock.RUnlock()

	builder := &strings.Builder{}
	builder.WriteByte('{')
	_, _ = fmt.Fprintf(builder, "min: %v, ", stat.min)
	_, _ = fmt.Fprintf(builder, "max: %v, ", stat.max)
	_, _ = fmt.Fprintf(builder, "mean: %v, ", mean(stat.total, stat.count))
	_, _ = fmt.Fprintf(builder, "total: %v, ", stat.count)
	_, _ = fmt.Fprintf(builder, "tps(1m): %v, ", stat.count1/interval1)
	_, _ = fmt.Fprintf(builder, "tps(5m): %v, ", stat.count5/interval5)
	_, _ = fmt.Fprintf(builder, "tps(10m): %v, ", stat.count10/interval10)
	_, _ = fmt.Fprintf(builder, "tps(60m): %v", stat.count60/interval60)
	builder.WriteByte('}')
	return builder.String()
}

// start 開始統計
func (this *Runtime) start() {
	go func() {
		timeout := time.After(time.Second)
		timeout1 := time.After(time.Second * interval1)
		timeout5 := time.After(time.Second * interval5)
		timeout10 := time.After(time.Second * interval10)
		timeout60 := time.After(time.Second * interval60)

		for {
			select {
			case <-timeout:
				this.lock.Lock()
				this.stat = this.curr
				this.lock.Unlock()

			case <-timeout1:
				this.lock.Lock()
				this.stat = this.curr
				this.curr.count1 = 0
				this.lock.Unlock()

			case <-timeout5:
				this.lock.Lock()
				this.stat = this.curr
				this.curr.count5 = 0
				this.lock.Unlock()

			case <-timeout10:
				this.lock.Lock()
				this.stat = this.curr
				this.curr.count10 = 0
				this.lock.Unlock()

			case <-timeout60:
				this.lock.Lock()
				this.stat = this.curr
				this.curr.count60 = 0
				this.lock.Unlock()

			default:
				if stop.Load() {
					return
				} // if
			} // select
		} // for
	}()
}

// mean 取得平均時間
func mean(total time.Duration, count int64) string {
	if count > 0 {
		return (total / time.Duration(count)).String()
	} else {
		return "n/a"
	} // if
}

// handle 取得http路由器, 會依照是否有認證資料來決定是否要添加帳號密碼認證
func handle(auth *Auth) http.Handler {
	handler := http.NewServeMux()

	if auth != nil {
		handler.Handle("/debug/vars", func(next http.Handler) http.Handler {
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
		handler.Handle("/debug/vars", expvar.Handler())
	} // if

	return handler
}

var stop atomic.Bool // 結束旗標
