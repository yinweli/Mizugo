package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find 搜尋行為
//
// 以匹配字串(Pattern)在次要資料庫中搜尋符合條件的索引, 並回傳結果列表(Data)
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MinorTable 與 MinorKey)
//   - 設定 Pattern: 不可為空字串, 將用於正則表達式匹配
//   - (可選)設定 Option: 用以調整 Mongo 查詢行為的選項列表
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - 建議搭配在欄位 MongoKey 上建立索引; 否則在資料量大時可能導致掃描成本偏高
//   - 查詢時會在合併 Option 後「強制加入 Projection」, 僅回傳欄位 MongoKey
//
// 工具函式 MinorIndex 建立專門用於次要資料庫的索引
//   - 自動以 Metaer.MinorTable 作為表格名稱
//   - 以 MongoKey 為索引欄位
//   - 排序固定為遞增(Order=1)
//   - 強制設為唯一索引(Unique=true)
//   - 索引名稱自動產生格式: "<MinorTable>_minor_index"
//
// 建議搭配 Find 行為使用, 可顯著提升 MongoKey 欄位的查詢效率
type Find struct {
	Behave                         // 行為物件
	Meta    Metaer                 // 元資料
	Pattern string                 // 匹配字串
	Option  []*options.FindOptions // 選項列表
	Data    []string               // 資料物件
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

	return nil
}

// Complete 完成處理
func (this *Find) Complete() error {
	table := this.Meta.MinorTable()
	result, err := this.Minor().Collection(table).Find(this.Ctx(), bson.M{MongoKey: bson.M{"$regex": this.Pattern}},
		append(this.Option, options.Find().SetProjection(bson.M{MongoKey: 1, "_id": 0}))...)

	if err != nil {
		return fmt.Errorf("find complete: %w: %v", err, this.Pattern)
	} // if

	defer func() {
		_ = result.Close(this.Ctx())
	}()
	this.Data = nil

	for result.Next(this.Ctx()) {
		temp := bson.M{}

		if err = result.Decode(&temp); err != nil {
			return fmt.Errorf("find complete: %w: %v", err, this.Pattern)
		} // if

		value, ok := temp[MongoKey].(string)

		if ok == false {
			return fmt.Errorf("find complete: field %s not found or not a string: %v", MongoKey, this.Pattern)
		} // if

		this.Data = append(this.Data, value)
	} // for

	if err = result.Err(); err != nil {
		return fmt.Errorf("find complete: %w: %v", err, this.Pattern)
	} // if

	return nil
}
