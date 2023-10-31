package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
)

// Del 刪除行為, 以索引字串與資料到主要/次要資料庫中刪除資料, 使用上有以下幾點須注意
//   - 需要事先建立好與 Metaer 介面符合的元資料結構, 並填寫到 Meta
//   - 執行前設定好 Key 並且不能為空字串
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

	if this.MinorEnable {
		if this.Meta.MinorTable() == "" {
			return fmt.Errorf("del prepare: table empty")
		} // if

		if this.Meta.MinorField() == "" {
			return fmt.Errorf("del prepare: field empty")
		} // if
	} // if

	return nil
}

// Complete 完成處理
func (this *Del) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("del complete: meta nil")
	} // if

	if this.MajorEnable {
		if _, err := this.cmd.Result(); err != nil {
			return fmt.Errorf("del complete: %w: %v", err, this.Key)
		} // if
	} // if

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		field := this.Meta.MinorField()
		filter := bson.D{{Key: field, Value: key}}

		if _, err := this.Minor().Table(table).DeleteOne(this.Ctx(), filter); err != nil {
			return fmt.Errorf("del complete: %w: %v", err, this.Key)
		} // if
	} // if

	return nil
}
