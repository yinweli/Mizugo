package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
)

// newMinor 建立次要資料庫, 並且連線到 MongoURI 指定的資料庫;
// 另外需要指定mongo資料庫名稱, 簡化後面取得執行器的流程, 但也因此限制次要資料庫不能在多個mongo資料庫間切換
func newMinor(ctx ctxs.Ctx, uri MongoURI, dbName string) (major *Minor, err error) {
	client, err := uri.Connect(ctx)

	if err != nil {
		return nil, fmt.Errorf("newMinor: %w", err)
	} // if

	if dbName == "" {
		return nil, fmt.Errorf("newMinor: dbName empty")
	} // if

	return &Minor{client: client, dbName: dbName}, nil
}

// Minor 次要資料庫, 內部用mongo實現的資料庫組件, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器, 實際上就是mongo集合
//   - 取得客戶端物件: 取得原生資料庫執行器, 可用來執行更細緻的命令
type Minor struct {
	client *mongo.Client // 客戶端物件
	dbName string        // 資料庫名稱
}

// MinorSubmit 資料庫執行器, 實際上就是mongo集合
type MinorSubmit = *mongo.Collection

// Submit 取得執行物件
func (this *Minor) Submit(tableName string) MinorSubmit {
	if this.client != nil {
		if database := this.client.Database(this.dbName); database != nil {
			return database.Collection(tableName)
		} // if
	} // if

	return nil
}

// Client 取得客戶端物件
func (this *Minor) Client() *mongo.Client {
	return this.client
}

// stop 停止資料庫
func (this *Minor) stop(ctx ctxs.Ctx) {
	if this.client != nil {
		_ = this.client.Disconnect(ctx.Ctx())
		this.client = nil
	} // if
}
