package redmos

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RankGet 排行榜取得名次行為
//
// 在次要資料庫中查詢, 取得指定索引在排行榜中的名次(Rank)
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MinorTable)
//   - 設定 Key: 玩家的索引值, 對應 MinorKey, 不可為空
//   - 設定 Ahead: 排序函式, 回傳「排在此索引前面」的 filter, 不可為 nil
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - 若索引不在榜上, Rank 為 0, 不回傳錯誤
//   - 「排在前面」的條件完全由呼叫端的 Ahead 函式決定, 組件本身不感知欄位語意
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
type RankGet[T any] struct {
	Behave                      // 行為物件
	Meta   Metaer               // 元資料
	Key    string               // 索引值(對應 MinorKey)
	Ahead  func(data *T) bson.M // 排序函式, 回傳「排在此索引前面」的 filter
	Rank   int64                // 結果名次; 0 表示索引不在榜上
}

// Prepare 前置處理
func (this *RankGet[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("rankget prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("rankget prepare: key empty")
	} // if

	if this.Ahead == nil {
		return fmt.Errorf("rankget prepare: ahead nil")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("rankget prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *RankGet[T]) Complete() error {
	key := this.Meta.MinorKey(this.Key)
	table := this.Meta.MinorTable()
	data := new(T)

	if err := this.Minor().Collection(table).FindOne(this.Ctx(), bson.M{MongoKey: key}).Decode(&MinorData[T]{D: data}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			this.Rank = 0
			return nil
		} // if

		return fmt.Errorf("rankget complete: %w", err)
	} // if

	count, err := this.Minor().Collection(table).CountDocuments(this.Ctx(), this.Ahead(data))

	if err != nil {
		return fmt.Errorf("rankget complete: %w", err)
	} // if

	this.Rank = count + 1
	return nil
}
