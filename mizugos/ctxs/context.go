package ctxs

import (
	"context"
	"time"
)

// Root 取得 Ctx 根物件, 當使用mizugo框架並且需要使用context時, 應該要由此函式的 Ctx 物件衍生
func Root() Ctx {
	return root
}

// RootCtx 方便取得 Ctx 根物件中的context.Context
func RootCtx() context.Context {
	return Root().Ctx()
}

// Ctx context資料
type Ctx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// Ctx 取得context物件
func (this Ctx) Ctx() context.Context {
	return this.ctx
}

// Done 檢查 Ctx 是否完成, 需搭配 select 語法或是 for range 等語法使用
func (this Ctx) Done() <-chan struct{} {
	return this.ctx.Done()
}

// Cancel 關閉由 Ctx 衍生出來的context, 也包含自身
func (this Ctx) Cancel() {
	if this.cancel != nil {
		this.cancel()
	} // if
}

// WithCancel 衍生包含取消功能的 Ctx
func (this Ctx) WithCancel() Ctx {
	result := Ctx{}
	result.ctx, result.cancel = context.WithCancel(this.ctx)
	return result
}

// WithTimeout 衍生包含超時功能的 Ctx
func (this Ctx) WithTimeout(duration time.Duration) Ctx {
	result := Ctx{}
	result.ctx, result.cancel = context.WithTimeout(this.ctx, duration)
	return result
}

// WithDeadline 衍生包含期限功能的 Ctx
func (this Ctx) WithDeadline(deadline time.Time) Ctx {
	result := Ctx{}
	result.ctx, result.cancel = context.WithDeadline(this.ctx, deadline)
	return result
}

func init() { //nolint:gochecknoinits
	root = Ctx{}
	root.ctx, root.cancel = context.WithCancel(context.Background())
}

var root Ctx // Ctx 根物件
