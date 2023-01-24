package miscs

import (
	"context"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/support/test_clientgo/internal/features"
)

// NewGenerator 建立連線產生器
func NewGenerator(max int, baseline, internal time.Duration, done func()) *Generator {
	return &Generator{
		max:      max,
		baseline: baseline,
		internal: internal,
		done:     done,
	}
}

// Generator 連線產生器
type Generator struct {
	max      int           // 最大連線數
	baseline time.Duration // 基準時間
	internal time.Duration // 間隔時間
	done     func()        // 完成物件
	cancel   func()        // 取消物件
	usage    usage         // 使用資料
}

// Start 啟動連線產生
func (this *Generator) Start() {
	mizugos.Poolmgr().Submit(func() {
		timeout := time.NewTicker(this.internal)
		ctx, cancel := context.WithCancel(contexts.Ctx())
		this.cancel = cancel

		for {
			select {
			case <-timeout.C:
				session := mizugos.Netmgr().Status().Session
				features.Connect.Set(int64(session))

				if this.max > session && this.usage.average() <= this.baseline {
					this.done()
				} // if

			case <-ctx.Done():
				timeout.Stop()
				return
			} // select
		} // for
	})
}

// Stop 停止連線產生
func (this *Generator) Stop() {
	if this.cancel != nil {
		this.cancel()
	} // if
}

// Report 回報時間
func (this *Generator) Report(duration time.Duration) {
	this.usage.add(duration)
}

// usage 使用資料
type usage struct {
	time  time.Duration // 使用時間
	count int           // 使用次數
	lock  sync.RWMutex  // 執行緒鎖
}

// add 新增使用時間與次數
func (this *usage) add(duration time.Duration) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.time += duration
	this.count++
}

// average 取得平均使用時間
func (this *usage) average() time.Duration {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if this.count > 0 {
		return this.time / time.Duration(this.count)
	} else {
		return 0
	} // if
}
