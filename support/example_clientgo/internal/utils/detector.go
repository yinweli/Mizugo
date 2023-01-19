package utils

import (
	"context"
	"time"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/support/example_clientgo/internal/features"
)

// Detector 連線檢測器
type Detector struct {
	notify chan any // 通知通道
	cancel func()   // 取消物件
}

// Start 啟動連線檢測
func (this *Detector) Start(count int, interval time.Duration, done func()) {
	mizugos.Poolmgr().Submit(func() {
		conn := func() {
			session := mizugos.Netmgr().Status().Session
			features.Connect.Set(int64(session))

			if session < count {
				done()
			} // if
		}
		timeout := time.NewTicker(interval)
		ctx, cancel := context.WithCancel(contexts.Ctx())
		this.notify = make(chan any, 1)
		this.cancel = cancel

		for {
			select {
			case <-this.notify:
				conn()

			case <-timeout.C:
				conn()

			case <-ctx.Done():
				timeout.Stop()
				return
			} // select
		} // for
	})
}

// Stop 停止連線檢測
func (this *Detector) Stop() {
	if this.cancel != nil {
		this.cancel()
	} // if
}

// Notice 通知連線變化
func (this *Detector) Notice() {
	if this.cancel != nil {
		this.notify <- nil
	} // if
}
