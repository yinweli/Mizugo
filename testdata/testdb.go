package testdata

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestDB 測試資料庫
type TestDB struct {
	key []string // 索引列表
}

// Reset 重置索引
func (this *TestDB) Reset() {
	this.key = []string{}
}

// Key 取得索引
func (this *TestDB) Key(key string) string {
	key = "test:" + key
	this.key = append(this.key, key)
	return key
}

// RedisClear 清除redis
func (this *TestDB) RedisClear(ctx context.Context, client redis.UniversalClient) {
	for _, itor := range this.key {
		client.Del(ctx, itor)
	} // for
}

// MongoClear 清除mongo, 直接刪除表格
func (this *TestDB) MongoClear(ctx context.Context, table *mongo.Collection) {
	_ = table.Drop(ctx)
}
