package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// Find 搜尋行為, 以匹配字串到次要資料庫中取得索引, 使用上有以下幾點須注意
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Pattern 並且不能為空字串
//   - 執行後可用 Data 來取得資料
//
// # Pattern匹配規則
//   - 使用正則表達式來寫匹配字串
type Find struct {
	Behave           // 行為物件
	Meta    Metaer   // 元資料
	Pattern string   // 匹配字串
	Data    []string // 資料物件
}

// Prepare 前置處理
func (this *Find) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("find prepare: meta nil")
	} // if

	if this.Pattern == "" {
		return fmt.Errorf("find prepare: pattern empty")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("find prepare: table empty")
	} // if

	if this.Meta.MinorField() == "" {
		return fmt.Errorf("find prepare: field empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *Find) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("find prepare: meta nil")
	} // if

	table := this.Meta.MinorTable()
	field := this.Meta.MinorField()
	filter := bson.M{field: bson.M{"$regex": this.Pattern}}
	result, err := this.Minor().Collection(table).Find(this.Ctx(), filter)

	if err != nil {
		return fmt.Errorf("find complete: %w: %v", err, this.Pattern)
	} // if

	defer func() {
		_ = result.Close(this.Ctx())
	}()
	data := []string{}

	for result.Next(this.Ctx()) {
		temp := bson.M{}

		if err = result.Decode(&temp); err != nil {
			return fmt.Errorf("find complete: %w: %v", err, this.Pattern)
		} // if

		value, ok := temp[field].(string)

		if ok == false {
			return fmt.Errorf("find complete: field %s not found or not a string: %v", field, this.Pattern)
		} // if

		data = append(data, value)
	} // for

	if err = result.Err(); err != nil {
		return fmt.Errorf("find complete: %w: %v", err, this.Pattern)
	} // if

	this.Data = data
	return nil
}
