package contexts

import (
	"context"
)

// Ctx 取得作為mizugo的context的根物件, 當使用mizugo框架並且需要使用context時, 應該要由此函式衍生context
func Ctx() context.Context {
	return ctx
}

// Cancel 關閉由 Ctx 衍生出來的context, 這是mizugo框架用來確保執行緒都被關閉的最後手段
func Cancel() {
	cancel()
}

func init() { //nolint
	ctx, cancel = context.WithCancel(context.Background())
}

var ctx context.Context       // ctx物件
var cancel context.CancelFunc // 取消物件
