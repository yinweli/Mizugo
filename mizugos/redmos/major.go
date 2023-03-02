package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
)

// newMajor 建立主要資料庫, 並且連線到 RedisURI 指定的資料庫
func newMajor(ctx ctxs.Ctx, uri RedisURI) (major *Major, err error) {
	client, err := uri.Connect(ctx)

	if err != nil {
		return nil, fmt.Errorf("newMajor: %w", err)
	} // if

	return &Major{client: client}, nil
}

// Major 主要資料庫, 內部用redis實現的資料庫組件, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器, 實際上就是redis管線
//   - 取得客戶端物件: 取得原生資料庫執行器, 可用來執行更細緻的命令
type Major struct {
	client redis.UniversalClient // 客戶端物件
}

// MajorSubmit 資料庫執行器, 實際上就是redis管線
type MajorSubmit = redis.Pipeliner

// Submit 取得執行物件
func (this *Major) Submit() MajorSubmit {
	if this.client != nil {
		return this.client.Pipeline()
	} // if

	return nil
}

// Client 取得客戶端物件
func (this *Major) Client() redis.UniversalClient {
	return this.client
}

// stop 停止資料庫
func (this *Major) stop() {
	if this.client != nil {
		_ = this.client.Close()
		this.client = nil
	} // if
}
