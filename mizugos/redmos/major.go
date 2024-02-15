package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
)

// newMajor 建立主要資料庫, 並且連線到 RedisURI 指定的資料庫
func newMajor(uri RedisURI) (major *Major, err error) {
	client, err := uri.Connect(ctxs.Get().Ctx())

	if err != nil {
		return nil, fmt.Errorf("newMajor: %w", err)
	} // if

	major = &Major{}
	major.client = client
	major.dbID = 0
	return major, nil
}

// Major 主要資料庫, 內部用redis實現的資料庫組件, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器, 實際上就是redis管線
//   - 取得客戶端物件: 取得原生資料庫執行器, 可用來執行更細緻的命令
type Major struct {
	client redis.UniversalClient // 客戶端物件
	dbID   int                   // 資料庫編號
}

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

// DBID 取得資料庫編號
func (this *Major) DBID() int {
	return this.dbID
}

// SwitchDB 切換資料庫, redis預設只能使用編號0~15的資料庫
func (this *Major) SwitchDB(dbID int) error {
	if this.client == nil {
		return fmt.Errorf("major switch: client nil")
	} // if

	if err := this.client.Do(ctxs.Get().Ctx(), "SELECT", dbID).Err(); err != nil {
		return fmt.Errorf("major switch: %w", err)
	} // if

	this.dbID = dbID
	return nil
}

// DropDB 清除資料庫
func (this *Major) DropDB() {
	if this.client != nil {
		this.client.FlushDB(ctxs.Get().Ctx())
	} // if
}

// stop 停止資料庫
func (this *Major) stop() {
	if this.client != nil {
		_ = this.client.Close()
		this.client = nil
	} // if
}

// MajorSubmit 資料庫執行器, 實際上就是redis管線
type MajorSubmit = redis.Pipeliner
