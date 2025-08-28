package redmos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetNX 設值(僅在不存在)行為
//
// 以索引鍵(Key)將資料寫入主要資料庫與/或次要資料庫, 但僅於主要資料庫(快取層)「鍵不存在」時才執行寫入
//
// 事前準備:
//   - 設定 MajorEnable / MinorEnable: 指示要作用的層
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(提供 MajorKey/MinorKey/MinorTable)
//   - 設定 Expire: 0 → 表示「不逾期」, RedisTTL → 表示「不更動逾期時間」, > 0 表示「設定新的逾期時間」
//   - 設定 Key: 不可為空字串
//   - 設定 Data: 不可為 nil, 且其成員需具備正確 bson 標籤(寫入次要資料庫時使用)
//
// 注意:
//   - 本行為包括儲存以及更新逾期時間
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
//   - 若 Data 實作 Saver 且 Saver.GetSave == false, 將直接略過寫入
//   - 僅在「主要資料庫鍵不存在」時才會寫入主要資料庫; 若鍵已存在, 行為會回傳錯誤並結束(不會寫入次要資料庫)
//   - 若只啟用次要資料庫, 則必定寫入次要資料庫
type SetNX[T any] struct {
	Behave                     // 行為物件
	MajorEnable bool           // 啟用主要資料庫
	MinorEnable bool           // 啟用次要資料庫
	Meta        Metaer         // 元資料
	Expire      time.Duration  // 逾期時間
	Key         string         // 索引值
	Data        *T             // 資料物件
	cmd         *redis.BoolCmd // 命令結果
}

// Prepare 前置處理
func (this *SetNX[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("setnx prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("setnx prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("setnx prepare: data nil")
	} // if

	if save, ok := any(this.Data).(Saver); ok && save.GetSave() == false {
		return nil
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		data, err := json.Marshal(this.Data)

		if err != nil {
			return fmt.Errorf("setnx prepare: %w: %v", err, this.Key)
		} // if

		this.cmd = this.Major().SetNX(this.Ctx(), key, data, this.Expire)
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("setnx prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *SetNX[T]) Complete() error {
	if save, ok := any(this.Data).(Saver); ok && save.GetSave() == false {
		return nil
	} // if

	if this.MajorEnable {
		result, err := this.cmd.Result()

		if err != nil {
			return fmt.Errorf("setnx complete: %w: %v", err, this.Key)
		} // if

		if result == false {
			return fmt.Errorf("setnx complete: save to redis failed: %v", this.Key)
		} // if
	} // if

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{MongoKey: key}).
			SetReplacement(&MinorData[T]{
				K: key,
				D: this.Data,
			}))
	} // if

	return nil
}
