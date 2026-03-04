package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexComplex 建立複合索引行為
//
// 在次要資料庫針對指定表格(Table)與多個欄位(Key)建立「複合索引」, 若索引已存在(以名稱比對)則不重複建立;
// 本行為僅影響次要資料庫的查詢效率, 不會更動資料內容, 且屬於可重入/冪等(idempotent)操作
//
// 事前準備:
//   - 設定 Name: 索引名稱(不可為空)
//   - 設定 Table: 目標表格名稱(不可為空)
//   - 設定 Key: 複合索引欄位列表(不可為空, Field 不可為空, Order 僅允許 1 或 -1)
//   - 設定 Unique: 是否唯一索引(true 則同一欄位值組合不可重複)
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - 若索引已存在同名索引, 將直接結束而不重建(以 Name 索引名稱做判斷)
//   - Unique=true 時, 若現有資料違反唯一性, Mongo 將拒絕建立索引並回傳錯誤
type IndexComplex struct {
	Behave             // 行為物件
	Name   string      // 索引名稱
	Table  string      // 表格名稱
	Key    []SortField // 複合索引欄位, 由外部指定
	Unique bool        // 是否唯一索引, 唯一索引的情況下, 索引值不允許重複
}

// Prepare 前置處理
func (this *IndexComplex) Prepare() error {
	if this.Name == "" {
		return fmt.Errorf("indexcomplex prepare: name empty")
	} // if

	if this.Table == "" {
		return fmt.Errorf("indexcomplex prepare: table empty")
	} // if

	if len(this.Key) == 0 {
		return fmt.Errorf("indexcomplex prepare: key empty")
	} // if

	for _, itor := range this.Key {
		if itor.Field == "" {
			return fmt.Errorf("indexcomplex prepare: key field empty")
		} // if

		if itor.Order != 1 && itor.Order != -1 {
			return fmt.Errorf("indexcomplex prepare: key order invalid")
		} // if
	} // for

	return nil
}

// Complete 完成處理
func (this *IndexComplex) Complete() error {
	collection := this.Minor().Collection(this.Table)
	result, err := collection.Indexes().List(this.Ctx())

	if err != nil {
		return fmt.Errorf("indexcomplex complete: %w: %v", err, this.Table)
	} // if

	defer func() {
		_ = result.Close(this.Ctx())
	}()

	for result.Next(this.Ctx()) {
		r := bson.M{}

		if result.Decode(&r) != nil {
			continue
		} // if

		name, ok := r["name"]

		if ok == false {
			continue
		} // if

		if name.(string) == this.Name {
			return nil // 如果索引已存在就結束了
		} // if
	} // for

	if err = result.Err(); err != nil {
		return fmt.Errorf("indexcomplex complete: %w: %v", err, this.Table)
	} // if

	key := bson.D{}

	for _, itor := range this.Key {
		key = append(key, bson.E{Key: itor.Field, Value: itor.Order})
	} // for

	if _, err = collection.Indexes().CreateOne(this.Ctx(), mongo.IndexModel{
		Keys:    key,
		Options: options.Index().SetName(this.Name).SetUnique(this.Unique),
	}); err != nil {
		return fmt.Errorf("indexcomplex complete: %w: %v", err, this.Table)
	} // if

	return nil
}
