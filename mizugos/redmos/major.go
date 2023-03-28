package redmos

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

// newMajor 建立主要資料庫, 並且連線到 RedisURI 指定的資料庫
func newMajor(context context.Context, uri RedisURI, record bool) (major *Major, err error) {
	client, err := uri.Connect(context)

	if err != nil {
		return nil, fmt.Errorf("newMajor: %w", err)
	} // if

	major = &Major{}
	major.client = client

	if record {
		list := &keylist{}
		client.AddHook(list)
		major.keylist = list
	} // if

	return major, nil
}

// Major 主要資料庫, 內部用redis實現的資料庫組件, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器, 實際上就是redis管線
//   - 取得客戶端物件: 取得原生資料庫執行器, 可用來執行更細緻的命令
type Major struct {
	client  redis.UniversalClient // 客戶端物件
	keylist *keylist              // 索引列表物件
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

// UsedKey 取得使用過的索引列表, 必須在建立時指定記錄旗標才有效, 否則就只會回傳空列表
func (this *Major) UsedKey() []string {
	if this.keylist != nil {
		return this.keylist.get()
	} // if

	return []string{}
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

// keylist 索引列表
type keylist struct {
	key []string // 索引列表
}

// DialHook redis hook
func (this *keylist) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook redis hook
func (this *keylist) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		this.add([]redis.Cmder{cmd})
		return next(ctx, cmd)
	}
}

// ProcessPipelineHook redis hook
func (this *keylist) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmd []redis.Cmder) error {
		this.add(cmd)
		return next(ctx, cmd)
	}
}

// add 增加索引
func (this *keylist) add(cmd []redis.Cmder) {
	for _, itor := range cmd {
		if arg := itor.Args(); len(arg) >= 2 {
			if value, ok := arg[1].(string); ok {
				this.key = append(this.key, value)
			} // if
		} // if
	} // for
}

// get 取得索引, 然後清空儲存的索引列表
func (this *keylist) get() []string {
	result := this.key
	this.key = []string{}
	return result
}
