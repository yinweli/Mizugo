package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Sample 取樣行為
//
// 使用 MongoDB $sample 聚合管道, 在次要資料庫中隨機取得 Count 筆資料並回傳(Data)
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MinorTable)
//   - 設定 Count: 取樣數量, 必須大於零
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - MongoDB $sample 在取樣數量小於集合 5% 時使用隨機排序; 超過則採兩階段演算法
//   - 若集合內文件數量少於 Count, 則 Data 筆數會少於 Count, 不回傳錯誤
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
type Sample[T any] struct {
	Behave        // 行為物件
	Meta   Metaer // 元資料
	Count  int    // 取樣數量
	Data   []*T   // 資料物件
}

// Prepare 前置處理
func (this *Sample[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("sample prepare: meta nil")
	} // if

	if this.Count <= 0 {
		return fmt.Errorf("sample prepare: count invalid")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("sample prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *Sample[T]) Complete() error {
	table := this.Meta.MinorTable()
	result, err := this.Minor().Collection(table).Aggregate(this.Ctx(), mongo.Pipeline{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: int64(this.Count)}}}},
	})

	if err != nil {
		return fmt.Errorf("sample complete: %w", err)
	} // if

	defer func() {
		_ = result.Close(this.Ctx())
	}()
	this.Data = nil

	for result.Next(this.Ctx()) {
		temp := new(T)

		if err = result.Decode(&MinorData[T]{
			D: temp,
		}); err != nil {
			return fmt.Errorf("sample complete: %w", err)
		} // if

		this.Data = append(this.Data, temp)
	} // for

	if err = result.Err(); err != nil {
		return fmt.Errorf("sample complete: %w", err)
	} // if

	return nil
}
