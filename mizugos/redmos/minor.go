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

	return &Minor{client: client, database: client.Database(dbName)}, nil
}

// Minor 次要資料庫, 內部用mongo實現的資料庫組件, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器, 實際上就是mongo集合
//   - 取得客戶端物件: 取得原生客戶端執行器, 可用來執行更細緻的命令
//   - 取得資料庫物件: 取得原生資料庫執行器, 可用來執行更細緻的命令
type Minor struct {
	client   *mongo.Client   // 客戶端物件
	database *mongo.Database // 資料庫物件
}

// Submit 取得執行物件
func (this *Minor) Submit() *MinorSubmit {
	if this.client != nil && this.database != nil {
		return &MinorSubmit{
			database: this.database,
		}
	} // if

	return nil
}

// Client 取得客戶端物件
func (this *Minor) Client() *mongo.Client {
	return this.client
}

// Database 取得資料庫物件
func (this *Minor) Database() *mongo.Database {
	return this.database
}

// stop 停止資料庫
func (this *Minor) stop(ctx ctxs.Ctx) {
	if this.client != nil {
		_ = this.client.Disconnect(ctx.Ctx())
		this.client = nil
		this.database = nil
	} // if
}

// MinorSubmit 資料庫執行器, 實際上就是mongo集合;
// 使用時必須先執行 Table 來初始化表格物件, 否則會造成錯誤
type MinorSubmit struct {
	database          *mongo.Database // 資料庫物件
	*mongo.Collection                 // 表格物件
}

// Table 初始化表格物件, 執行之後才能開始用表格操作
func (this *MinorSubmit) Table(tableName string) *MinorSubmit {
	this.Collection = this.database.Collection(tableName)
	return this
}
