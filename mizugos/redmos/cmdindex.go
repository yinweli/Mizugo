package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Index 建立索引行為, 到次要資料庫中建立索引, 使用上有以下幾點須注意
//   - 只有次要資料庫的操作會被索引影響, 當查詢的欄位符合索引時會自動生效
//   - 執行前設定好 Name, Table, Field, Order, Unique 並且要符合規範
type Index struct {
	Behave        // 行為物件
	Name   string // 索引名稱
	Table  string // 表格名稱
	Field  string // 欄位名稱
	Order  int    // 排序方向, 1表示順序, -1表示逆序
	Unique bool   // 是否唯一索引, 唯一索引的情況下, 索引值不允許重複
}

// Prepare 前置處理
func (this *Index) Prepare() error {
	if this.Name == "" {
		return fmt.Errorf("index prepare: name empty")
	} // if

	if this.Table == "" {
		return fmt.Errorf("index prepare: table empty")
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
	collection := this.Minor().Collection(this.Table)
	index, err := collection.Indexes().List(this.Ctx())

	if err != nil {
		return fmt.Errorf("index complete: %w: %v(%v)", err, this.Table, this.Field)
	} // if

	for index.Next(this.Ctx()) {
		x := bson.M{}

		if index.Decode(&x) != nil {
			continue
		} // if

		name, ok := x["name"]

		if ok == false {
			continue
		} // if

		if name.(string) == this.Name {
			return nil // 如果索引已存在就結束了
		} // if
	} // for

	if _, err = collection.Indexes().CreateOne(this.Ctx(), mongo.IndexModel{
		Keys:    bson.M{this.Field: this.Order},
		Options: options.Index().SetName(this.Name).SetUnique(this.Unique),
	}); err != nil {
		return fmt.Errorf("index complete: %w: %v(%v)", err, this.Table, this.Field)
	} // if

	return nil
}
