package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RankSet 排行榜設定分數行為
//
// 在次要資料庫中以條件式寫入排行榜資料:
//   - 文件不存在時插入($setOnInsert)
//   - 文件存在且符合 Filter 條件時才替換(ReplaceOne)
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MinorKey/MinorTable)
//   - 設定 Key: 不可為空字串
//   - 設定 Filter: 描述「現有文件需被更新」的條件, 不需包含索引欄位(會自動合併)
//   - 設定 Data: 不可為 nil, 且其成員需具備正確 bson 標籤(寫入次要資料庫時使用)
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
type RankSet[T any] struct {
	Behave        // 行為物件
	Meta   Metaer // 元資料
	Key    string // 索引值
	Filter bson.M // 條件篩選: 描述「現有文件需被更新」的情況
	Data   *T     // 資料物件
}

// Prepare 前置處理
func (this *RankSet[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("rankset prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("rankset prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("rankset prepare: data nil")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("rankset prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *RankSet[T]) Complete() error {
	key := this.Meta.MinorKey(this.Key)
	table := this.Meta.MinorTable()
	raw, err := bson.Marshal(&MinorData[T]{K: key, D: this.Data})

	if err != nil {
		return fmt.Errorf("rankset complete: %w: %v", err, this.Key)
	} // if

	doc := bson.M{}

	if err = bson.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("rankset complete: %w: %v", err, this.Key)
	} // if

	filter := bson.M{MongoKey: key}

	for k, v := range this.Filter {
		filter[k] = v
	} // for

	this.Minor().Operate(table, mongo.NewUpdateOneModel().
		SetUpsert(true).
		SetFilter(bson.M{MongoKey: key}).
		SetUpdate(bson.M{"$setOnInsert": doc}))
	this.Minor().Operate(table, mongo.NewReplaceOneModel().
		SetFilter(filter).
		SetReplacement(&MinorData[T]{
			K: key,
			D: this.Data,
		}))
	return nil
}
