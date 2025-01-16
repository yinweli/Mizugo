package redmos

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// newMinor 建立次要資料庫, 並且連線到 MongoURI 指定的資料庫
func newMinor(uri MongoURI, dbName string) (major *Minor, err error) {
	client, err := uri.Connect(context.Background())

	if err != nil {
		return nil, fmt.Errorf("newMinor: %w", err)
	} // if

	if dbName == "" {
		return nil, fmt.Errorf("newMinor: dbName empty")
	} // if

	minor := &Minor{}
	minor.client = client
	minor.database = client.Database(dbName)
	return minor, nil
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
			operate:  map[string][]mongo.WriteModel{},
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

// SwitchDB 切換資料庫
func (this *Minor) SwitchDB(dbName string) error {
	if this.client == nil {
		return fmt.Errorf("minor switch: client nil")
	} // if

	if dbName == "" {
		return fmt.Errorf("minor switch: dbName empty")
	} // if

	this.database = this.client.Database(dbName)
	return nil
}

// DropDB 清除資料庫
func (this *Minor) DropDB() {
	if this.client != nil {
		_ = this.database.Drop(context.Background())
	} // if
}

// stop 停止資料庫
func (this *Minor) stop() {
	if this.client != nil {
		_ = this.client.Disconnect(context.Background())
		this.client = nil
		this.database = nil
	} // if
}

// MinorSubmit 資料庫執行器
type MinorSubmit struct {
	database *mongo.Database               // 資料庫物件
	operate  map[string][]mongo.WriteModel // 批量操作列表
}

// Collection 取得表格物件
func (this *MinorSubmit) Collection(table string) *mongo.Collection {
	return this.database.Collection(table)
}

// Operate 新增批量操作
func (this *MinorSubmit) Operate(table string, operate mongo.WriteModel) *MinorSubmit {
	this.operate[table] = append(this.operate[table], operate)
	return this
}

// Exec 執行批量操作
func (this *MinorSubmit) Exec(ctx context.Context) error {
	for table, itor := range this.operate {
		if _, err := this.database.Collection(table).BulkWrite(ctx, itor); err != nil {
			return fmt.Errorf("minorSubmit exec: %w", err)
		} // if
	} // for

	this.operate = map[string][]mongo.WriteModel{}
	return nil
}
