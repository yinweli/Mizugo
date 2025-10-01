package redmos

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// newMajor 建立主要資料庫, 並依據傳入的 RedisURI 立刻進行連線
func newMajor(uri RedisURI) (major *Major, err error) {
	client, err := uri.Connect(context.Background())

	if err != nil {
		return nil, fmt.Errorf("newMajor: %w", err)
	} // if

	return &Major{
		client: client,
		uri:    uri,
	}, nil
}

// Major 主要資料庫(Redis)
//
// 以 Redis 為基礎的資料庫封裝, 提供兩種使用模式:
//   - Submit: 取得 redis.Pipeliner, 適合批次/管線化命令
//   - Client: 取得原生 redis.UniversalClient, 適合需要原生 API 的場景
//
// Major 不是執行緒安全, 若跨 goroutine 共用, 請由上層管理器(Redmomgr)負責同步保護
type Major struct {
	client redis.UniversalClient // 客戶端物件
	uri    RedisURI              // 連接字串
}

// MajorSubmit Pipeline 執行器, 實際上是 redis.Pipeliner
//
// 用法:
//
//	pipe := major.Submit()
//	pipe.Set(ctx, "k", "v", 0)
//	_, err := pipe.Exec(ctx)
type MajorSubmit = redis.Pipeliner

// Submit 取得 Pipeline 執行器
func (this *Major) Submit() MajorSubmit {
	if this.client != nil {
		return this.client.Pipeline()
	} // if

	return nil
}

// Client 取得 Redis 客戶端
func (this *Major) Client() redis.UniversalClient {
	return this.client
}

// SwitchDB 切換 Redis DB
//
// 以目前 RedisURI 為基礎新增 'dbid=N', 並重新連線; 成功後會關閉舊 client, 再以新 client 取代
func (this *Major) SwitchDB(dbID int) error {
	if this.client == nil {
		return fmt.Errorf("major switch: client nil")
	} // if

	client, err := this.uri.add(fmt.Sprintf("dbid=%v", dbID)).Connect(context.Background())

	if err != nil {
		return fmt.Errorf("major switch: %w", err)
	} // if

	_ = this.client.Close()
	this.client = client
	return nil
}

// DropDB 清除資料庫
func (this *Major) DropDB() {
	if this.client != nil {
		this.client.FlushDB(context.Background())
	} // if
}

// stop 關閉並釋放資料庫
func (this *Major) stop() {
	if this.client != nil {
		_ = this.client.Close()
		this.client = nil
	} // if
}
