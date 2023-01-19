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
	cancel func() // 取消物件
}

// Start 啟動連線檢測
func (this *Detector) Start(total, batch int, interval time.Duration, done func()) {
	mizugos.Poolmgr().Submit(func() {
		timeout := time.NewTicker(interval)
		ctx, cancel := context.WithCancel(contexts.Ctx())
		this.cancel = cancel

		for {
			select {
			case <-timeout.C:
				this.connect(total, batch, done)

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

// connect 連線處理
func (this *Detector) connect(total, batch int, done func()) {
	session := mizugos.Netmgr().Status().Session
	features.Connect.Set(int64(session))

	if total <= session {
		return
	} // if

	count := total - session

	if count > batch {
		count = batch
	} // if

	for batch > 0 {
		done()
		batch--
	} // if
}
