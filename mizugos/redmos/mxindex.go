package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Index 建立索引行為, 到次要資料庫中建立索引, 使用上有以下幾點須注意
//   - 只有次要資料庫的操作會被索引影響, 當查詢的欄位符合索引時會自動生效
//   - 執行前必須設定好 Table, Field 並且不能為空字串
//   - 執行前設定好 Order, Unique 並且要符合規範
//   - 在內部執行過程中, 索引名稱, 索引欄位會被轉為小寫
type Index struct {
	Behave
	Table  string // 表格名稱
	Field  string // 索引名稱
	Order  int    // 排序方向, 1表示順序, -1表示逆序
	Unique bool   // 是否唯一索引, 唯一索引的情況下, 索引值不允許重複
}

// Prepare 前置處理
func (this *Index) Prepare() error {
	if this.Table == "" {
		return fmt.Errorf("index prepare: type empty")
	} // if

	if this.Field == "" {
		return fmt.Errorf("index prepare: field empty")
	} // if

	if this.Order != 1 && this.Order != -1 {
		return fmt.Errorf("index prepare: order invalid")
	} // if

	return nil
}

// Complete 完成處理
func (this *Index) Complete() error {
	model := mongo.IndexModel{
		Keys:    bson.M{this.Field: this.Order},
		Options: options.Index().SetUnique(this.Unique),
	}

	if _, err := this.Minor().Table(this.Table).Indexes().CreateOne(this.Ctx(), model); err != nil {
		return fmt.Errorf("index complete: %w", err)
	} // if

	return nil
}
