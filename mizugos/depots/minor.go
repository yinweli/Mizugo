package depots

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// newMinor 建立次要資料庫, 並且連線到 MongoURI 指定的資料庫
func newMinor(ctx context.Context, uri MongoURI) (major *Minor, err error) {
	client, err := uri.Connect(ctx)

	if err != nil {
		return nil, fmt.Errorf("newMinor: %w", err)
	} // if

	return &Minor{client: client}, nil
}

// Minor 次要資料庫, 內部用mongo實現的資料庫組件, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器, 實際上就是mongo集合
//   - 取得客戶端物件: 取得原生資料庫執行器, 可用來執行更細緻的命令
type Minor struct {
	client *mongo.Client // 客戶端物件
}

// MinorRunner 資料庫執行器, 實際上就是mongo集合
type MinorRunner = *mongo.Collection

// Runner 取得執行物件
func (this *Minor) Runner(dbName, tableName string) MinorRunner {
	if this.client != nil {
		if database := this.client.Database(dbName); database != nil {
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
func (this *Minor) stop(ctx context.Context) {
	if this.client != nil {
		_ = this.client.Disconnect(ctx)
		this.client = nil
	} // if
}
