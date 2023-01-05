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
	curr   runtime      // 當前資料
	last   runtime      // 上次資料
	lock   sync.RWMutex // 執行緒鎖
}

// runtime 執行資料
type runtime struct {
	time    time.Duration // 總執行時間
	max     time.Duration // 最大執行時間
	count   int64         // 總執行次數
	count1  int64         // 每分鐘執行次數
	count5  int64         // 每5分鐘執行次數
	count10 int64         // 每10分鐘執行次數
	count60 int64         // 每60分鐘執行次數
}

// Rec 執行紀錄, 使用時把回傳的函式指標記錄下來, 直到執行區間結束再執行
func (this *Runtime) Rec() func() {
	start := time.Now()
	return func() {
		delta := time.Since(start)

		this.lock.Lock()
		defer this.lock.Unlock()

		this.curr.time += delta

		if this.curr.max < delta {
			this.curr.max = delta
		} // if

		this.curr.count++
		this.curr.count1++
		this.curr.count5++
		this.curr.count10++
		this.curr.count60++
	}
}

// String 取得統計字串
func (this *Runtime) String() string {
	this.lock.RLock()
	stat := this.last
	this.lock.RUnlock()

	mean := "n/a"

	if stat.count > 0 {
		mean = (stat.time / time.Duration(stat.count)).String()
	} // if

	builder := &strings.Builder{}
	builder.WriteByte('{')
	_, _ = fmt.Fprintf(builder, "\"time\": \"%v\", ", stat.time)
	_, _ = fmt.Fprintf(builder, "\"max\": \"%v\", ", stat.max)
	_, _ = fmt.Fprintf(builder, "\"mean\": \"%v\", ", mean)
	_, _ = fmt.Fprintf(builder, "\"count\": %v, ", stat.count)
	_, _ = fmt.Fprintf(builder, "\"count(1m)\": %v, ", stat.count1)
	_, _ = fmt.Fprintf(builder, "\"count(5m)\": %v, ", stat.count5)
	_, _ = fmt.Fprintf(builder, "\"count(10m)\": %v, ", stat.count10)
	_, _ = fmt.Fprintf(builder, "\"count(60m)\": %v", stat.count60)
	builder.WriteByte('}')
	return builder.String()
}

// start 開始統計
func (this *Runtime) start() {
	go func() {
		timeout := time.NewTicker(time.Second)
		timeout1 := time.NewTicker(time.Second * interval1)
		timeout5 := time.NewTicker(time.Second * interval5)
		timeout10 := time.NewTicker(time.Second * interval10)
		timeout60 := time.NewTicker(time.Second * interval60)

		for {
			select {
			case <-timeout.C:
				this.lock.Lock()
				this.last.time = this.curr.time
				this.last.max = this.curr.max
				this.last.count = this.curr.count
				this.lock.Unlock()

			case <-timeout1.C:
				this.lock.Lock()
				this.last.count1 = this.curr.count1
				this.curr.count1 = 0
				this.lock.Unlock()

			case <-timeout5.C:
				this.lock.Lock()
				this.last.count5 = this.curr.count5
				this.curr.count5 = 0
				this.lock.Unlock()

			case <-timeout10.C:
				this.lock.Lock()
				this.last.count10 = this.curr.count10
				this.curr.count10 = 0
				this.lock.Unlock()

			case <-timeout60.C:
				this.lock.Lock()
				this.last.count60 = this.curr.count60
				this.curr.count60 = 0
				this.lock.Unlock()

			default:
				if this.finish() {
					timeout.Stop()
					timeout1.Stop()
					timeout5.Stop()
					timeout10.Stop()
					timeout60.Stop()
					return
				} // if
			} // select
		} // for
	}()
}
