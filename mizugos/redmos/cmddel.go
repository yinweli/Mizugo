package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Del 刪除行為
//
// 以索引鍵(Key)刪除主要資料庫與/或次要資料庫的對應紀錄
//
// 事前準備:
//   - 設定 MajorEnable / MinorEnable: 指示要作用的層
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(提供 MajorKey/MinorKey/MinorTable)
//   - 設定 Key: 不可為空字串
type Del struct {
	Behave                    // 行為物件
	MajorEnable bool          // 啟用主要資料庫
	MinorEnable bool          // 啟用次要資料庫
	Meta        Metaer        // 元資料
	Key         string        // 索引值
	cmd         *redis.IntCmd // 命令結果
}

// Prepare 前置處理
func (this *Del) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("del prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("del prepare: key empty")
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		this.cmd = this.Major().Del(this.Ctx(), key)
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("del prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *Del) Complete() error {
	if this.MajorEnable {
		if _, err := this.cmd.Result(); err != nil {
			return fmt.Errorf("del complete: %w: %v", err, this.Key)
		} // if
	} // if

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewDeleteOneModel().
			SetFilter(bson.M{MongoKey: key}))
	} // if

	return nil
}
