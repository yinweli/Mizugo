package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// Runtime 執行統計
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

// Add 新增紀錄
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

// Rec 執行紀錄, 使用時把回傳的函式指標記錄下來, 直到執行區間結束再執行
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

	return utils.ExpvarStr([]utils.ExpvarStat{
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

// start 開始統計
func (this *Runtime) start(ctx context.Context) {
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
				this.last.timeMax = this.curr.timeMax
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

			case <-ctx.Done():
				timeout.Stop()
				timeout1.Stop()
				timeout5.Stop()
				timeout10.Stop()
				timeout60.Stop()
				return
			} // select
		} // for
	}()
}
