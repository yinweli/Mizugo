package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Index 建立索引行為, 到次要資料庫中建立索引, 使用上有以下幾點須注意
//   - 只有次要資料庫的操作會被索引影響, 當查詢的欄位符合索引時會自動生效
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Name, Order, Unique 並且要符合規範
type Index struct {
	Behave        // 行為物件
	Meta   Metaer // 元資料
	Name   string // 索引名稱
	Order  int    // 排序方向, 1表示順序, -1表示逆序
	Unique bool   // 是否唯一索引, 唯一索引的情況下, 索引值不允許重複
}

// Prepare 前置處理
func (this *Index) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("index prepare: meta nil")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("index prepare: table empty")
	} // if

	if this.Meta.MinorField() == "" {
		return fmt.Errorf("index prepare: field empty")
	} // if

	if this.Name == "" {
		return fmt.Errorf("index prepare: name empty")
	} // if

	if this.Order != 1 && this.Order != -1 {
		return fmt.Errorf("index prepare: order invalid")
	} // if

	return nil
}

// Complete 完成處理
func (this *Index) Complete() error {
	table := this.Meta.MinorTable()
	field := this.Meta.MinorField()
	collection := this.Minor().Collection(table)
	index, err := collection.Indexes().List(this.Ctx())

	if err != nil {
		return fmt.Errorf("index complete: %w: %v(%v)", err, table, field)
	} // if

	for index.Next(this.Ctx()) {
		info := bson.M{}

		if index.Decode(&info) != nil {
			continue
		} // if

		name, ok := info["name"]

		if ok == false {
			continue
		} // if

		if name.(string) == this.Name {
			return nil // 如果索引已存在就結束了
		} // if
	} // for

	model := mongo.IndexModel{
		Keys:    bson.M{field: this.Order},
		Options: options.Index().SetName(this.Name).SetUnique(this.Unique),
	}

	if _, err = collection.Indexes().CreateOne(this.Ctx(), model); err != nil {
		return fmt.Errorf("index complete: %w: %v(%v)", err, table, field)
	} // if

	return nil
}
