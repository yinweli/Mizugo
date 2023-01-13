package contexts

import (
	"context"
)

// context, 這裡提供了mizugo中使用到context功能的根物件
// 當執行Cancel函式時會把從這裡(也就是Ctx函式)衍生出來的執行緒都關閉
// 這是mizugo用來確保執行緒都被關閉的最後手段

// Ctx 取得ctx物件
func Ctx() context.Context {
	return ctx
}

// Cancel 執行取消
func Cancel() {
	cancel()
}

func init() { //nolint
	ctx, cancel = context.WithCancel(context.Background())
}

var ctx context.Context       // ctx物件
var cancel context.CancelFunc // 取消物件
