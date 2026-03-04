package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RankList 排行榜列表行為
//
// 在次要資料庫中查詢, 依指定排序規則取得前 Limit 筆資料並回傳(Data)
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MinorTable)
//   - 設定 Limit: 取得筆數, 必須大於零
//   - 設定 Sort: 排序規則(Field 不可為空, Order 僅允許 1 或 -1)
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - 若集合內文件數量少於 Limit, 則 Data 筆數會少於 Limit, 不回傳錯誤
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
type RankList[T any] struct {
	Behave             // 行為物件
	Meta   Metaer      // 元資料
	Limit  int64       // 取得筆數, 必須大於零
	Sort   []SortField // 排序規則, 由外部指定, 不可為空
	Data   []*T        // 資料物件
}

// Prepare 前置處理
func (this *RankList[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("ranklist prepare: meta nil")
	} // if

	if this.Limit <= 0 {
		return fmt.Errorf("ranklist prepare: limit invalid")
	} // if

	if len(this.Sort) == 0 {
		return fmt.Errorf("ranklist prepare: sort empty")
	} // if

	for _, s := range this.Sort {
		if s.Field == "" {
			return fmt.Errorf("ranklist prepare: sort field empty")
		} // if

		if s.Order != 1 && s.Order != -1 {
			return fmt.Errorf("ranklist prepare: sort order invalid")
		} // if
	} // for

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("ranklist prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *RankList[T]) Complete() error {
	sort := bson.D{}

	for _, s := range this.Sort {
		sort = append(sort, bson.E{Key: s.Field, Value: s.Order})
	} // for

	table := this.Meta.MinorTable()
	option := options.Find().SetSort(sort).SetLimit(this.Limit)
	result, err := this.Minor().Collection(table).Find(this.Ctx(), bson.D{}, option)

	if err != nil {
		return fmt.Errorf("ranklist complete: %w", err)
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
			return fmt.Errorf("ranklist complete: %w", err)
		} // if

		this.Data = append(this.Data, temp)
	} // for

	if err = result.Err(); err != nil {
		return fmt.Errorf("ranklist complete: %w", err)
	} // if

	return nil
}
