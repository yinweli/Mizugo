package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get 取值行為
//
// 以索引鍵(Key)讀取主要資料庫與/或次要資料庫中的資料, 僅查詢不寫入;
// 若主要資料庫(快取層)命中則直接回傳, 否則回落至次要資料庫(持久層)查詢;
// 讀取主要資料庫時, 不會改變該鍵的逾期時間(TTL)
//
// 事前準備:
//   - 設定 MajorEnable / MinorEnable: 指示要作用的層
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(提供 MajorKey/MinorKey/MinorTable)
//   - 設定 Key: 不可為空字串
//   - (可選)設定 Data: 若為 nil, 執行時會自動建立 *T
//
// 注意:
//   - 本行為僅讀取, 不會刷新主要資料庫 TTL
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
//   - 若同時啟用 Major 與 Minor, 會先嘗試 Major; 命中即結束, 不再查 Minor
type Get[T any] struct {
	Behave                       // 行為物件
	MajorEnable bool             // 啟用主要資料庫
	MinorEnable bool             // 啟用次要資料庫
	Meta        Metaer           // 元資料
	Key         string           // 索引值
	Data        *T               // 資料物件
	cmd         *redis.StringCmd // 命令結果
}

// Prepare 前置處理
func (this *Get[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("get prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("get prepare: key empty")
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		this.cmd = this.Major().Get(this.Ctx(), key)
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("get prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *Get[T]) Complete() error {
	if this.MajorEnable {
		result, err := this.cmd.Result()

		if err != nil && errors.Is(err, redis.Nil) == false {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		if result != RedisNil {
			if this.Data == nil {
				this.Data = new(T)
			} // if

			if err = json.Unmarshal([]byte(result), this.Data); err != nil {
				return fmt.Errorf("get complete: %w: %v", err, this.Key)
			} // if

			return nil // 如果主要資料庫讀取成功, 就不必執行次要資料庫了
		} // if
	} // if

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		result := this.Minor().Collection(table).FindOne(this.Ctx(), bson.M{MongoKey: key})
		err := result.Err()
		empty := errors.Is(err, mongo.ErrNoDocuments)

		if err != nil && empty == false {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		if empty == false {
			if this.Data == nil {
				this.Data = new(T)
			} // if

			if err = result.Decode(&MinorData[T]{
				D: this.Data,
			}); err != nil {
				return fmt.Errorf("get complete: %w: %v", err, this.Key)
			} // if
		} // if
	} // if

	return nil
}
