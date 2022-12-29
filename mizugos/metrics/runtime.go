package metrics

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Runtime 執行統計
type Runtime struct {
	finish func() bool  // 終止檢驗函式
	stat   runtime      // 統計資料
	curr   runtime      // 當前資料
	lock   sync.RWMutex // 執行緒鎖
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

	mean := "n/a"

	if stat.count > 0 {
		mean = (stat.total / time.Duration(stat.count)).String()
	} // if

	builder := &strings.Builder{}
	builder.WriteByte('{')
	_, _ = fmt.Fprintf(builder, "min: %v, ", stat.min)
	_, _ = fmt.Fprintf(builder, "max: %v, ", stat.max)
	_, _ = fmt.Fprintf(builder, "mean: %v, ", mean)
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
				if this.finish() {
					return
				} // if
			} // select
		} // for
	}()
}
