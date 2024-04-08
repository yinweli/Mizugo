package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/mizugo/helps"
	"github.com/yinweli/Mizugo/mizugo/pools"
)

// Runtime 執行度量, 記錄了以下數據
//   - 總執行時間
//   - 最大執行時間
//   - 總執行次數
//   - 每分鐘執行次數
//   - 每5分鐘執行次數
//   - 每10分鐘執行次數
//   - 每60分鐘執行次數
type Runtime struct {
	curr runtime      // 當前資料
	last runtime      // 上次資料
	lock sync.RWMutex // 執行緒鎖
}

// runtime 執行資料
type runtime struct {
	time    time.Duration // 總執行時間
	timeMax time.Duration // 最大執行時間
	count   int64         // 總執行次數
	count1  int64         // 每分鐘執行次數
	count5  int64         // 每5分鐘執行次數
	count10 int64         // 每10分鐘執行次數
	count60 int64         // 每60分鐘執行次數
}

// Add 新增記錄
func (this *Runtime) Add(delta time.Duration) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.curr.time += delta

	if this.curr.timeMax < delta {
		this.curr.timeMax = delta
	} // if

	this.curr.count++
	this.curr.count1++
	this.curr.count5++
	this.curr.count10++
	this.curr.count60++
}

// Rec 執行記錄, 使用時把回傳的函式指標記錄下來, 直到執行區間結束再執行
func (this *Runtime) Rec() func() {
	start := time.Now()
	return func() {
		this.Add(time.Since(start))
	}
}

// String 取得字串
func (this *Runtime) String() string {
	this.lock.RLock()
	stat := this.last
	this.lock.RUnlock()

	timeAvg := "n/a"

	if stat.count > 0 {
		timeAvg = (stat.time / time.Duration(stat.count)).String()
	} // if

	return helps.ExpvarStr([]helps.ExpvarStat{
		{Name: "time", Data: stat.time},
		{Name: "time(max)", Data: stat.timeMax},
		{Name: "time(avg)", Data: timeAvg},
		{Name: "count", Data: stat.count},
		{Name: "count(1m)", Data: stat.count1},
		{Name: "count(5m)", Data: stat.count5},
		{Name: "count(10m)", Data: stat.count10},
		{Name: "count(60m)", Data: stat.count60},
	})
}

// start 開始度量
func (this *Runtime) start(ctx context.Context) {
	pools.DefaultPool.Submit(func() {
		ticker := time.NewTicker(time.Second)
		ticker1 := time.NewTicker(time.Second * interval1)
		ticker5 := time.NewTicker(time.Second * interval5)
		ticker10 := time.NewTicker(time.Second * interval10)
		ticker60 := time.NewTicker(time.Second * interval60)

		for {
			select {
			case <-ticker.C:
				this.lock.Lock()
				this.last.time = this.curr.time
				this.last.timeMax = this.curr.timeMax
				this.last.count = this.curr.count
				this.lock.Unlock()

			case <-ticker1.C:
				this.lock.Lock()
				this.last.count1 = this.curr.count1
				this.curr.count1 = 0
				this.lock.Unlock()

			case <-ticker5.C:
				this.lock.Lock()
				this.last.count5 = this.curr.count5
				this.curr.count5 = 0
				this.lock.Unlock()

			case <-ticker10.C:
				this.lock.Lock()
				this.last.count10 = this.curr.count10
				this.curr.count10 = 0
				this.lock.Unlock()

			case <-ticker60.C:
				this.lock.Lock()
				this.last.count60 = this.curr.count60
				this.curr.count60 = 0
				this.lock.Unlock()

			case <-ctx.Done():
				ticker.Stop()
				ticker1.Stop()
				ticker5.Stop()
				ticker10.Stop()
				ticker60.Stop()
				return
			} // select
		} // for
	})
}
