package contexts

import (
	"context"
)

// Ctx 取得ctx物件
func Ctx() context.Context {
	return ctx
}

// Cancel 執行取消
func Cancel() {
	cancel()
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

var ctx context.Context       // ctx物件
var cancel context.CancelFunc // 取消物件
