package redmos

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Index 建立索引行為
//
// 在次要資料庫針對指定表格(Table)與欄位(Field)建立「單欄位索引」, 若索引已存在(以名稱比對)則不重複建立;
// 本行為僅影響次要資料庫的查詢效率, 不會更動資料內容, 且屬於可重入/冪等(idempotent)操作
//
// 事前準備:
//   - 設定 Name: 索引名稱(不可為空)
//   - 設定 Table: 目標表格名稱(不可為空)
//   - 設定 Sort: 排序欄位(Field 不可為空, Order 僅允許 1 或 -1)
//   - 設定 Unique: 是否唯一索引(true 則同一欄位值不可重複)
//
// 注意:
//   - 本行為僅使用次要資料庫, 主要資料庫不參與
//   - 僅建立「單欄位索引」; 複合索引(多欄位)請改用 IndexComplex
//   - 若欄位已存在同名索引, 將直接結束而不重建(以 Name 索引名稱做判斷)
//   - Unique=true 時, 若現有資料違反唯一性, Mongo 將拒絕建立索引並回傳錯誤
//   - 索引能顯著改善查詢效率, 但也會增加寫入成本與存儲空間; 請評估實際讀寫比例再建立
//   - 在高資料量表上初次建立索引可能耗時; 建議於低峰期或維護時段操作
//
// 工具函式 MinorIndex 建立專門用於次要資料庫的索引
//   - 自動以 Metaer.MinorTable 作為表格名稱
//   - 以 MongoKey 為索引欄位, 排序固定為遞增(Order=1)
//   - 強制設為唯一索引(Unique=true)
//   - 索引名稱自動產生格式: "<MinorTable>_minor_index"
type Index struct {
	Behave           // 行為物件
	Name   string    // 索引名稱
	Table  string    // 表格名稱
	Sort   SortField // 排序欄位
	Unique bool      // 是否唯一索引, 唯一索引的情況下, 索引值不允許重複
}

// Prepare 前置處理
func (this *Index) Prepare() error {
	if this.Name == "" {
		return fmt.Errorf("index prepare: name empty")
	} // if

	if this.Table == "" {
		return fmt.Errorf("index prepare: table empty")
	} // if

	if this.Sort.Field == "" {
		return fmt.Errorf("index prepare: field empty")
	} // if

	if this.Sort.Order != 1 && this.Sort.Order != -1 {
		return fmt.Errorf("index prepare: order invalid")
	} // if

	return nil
}

// Complete 完成處理
func (this *Index) Complete() error {
	collection := this.Minor().Collection(this.Table)
	result, err := collection.Indexes().List(this.Ctx())

	if err != nil {
		return fmt.Errorf("index complete: %w: %v(%v)", err, this.Table, this.Sort.Field)
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
		return fmt.Errorf("index complete: %w: %v(%v)", err, this.Table, this.Sort.Field)
	} // if

	if _, err = collection.Indexes().CreateOne(this.Ctx(), mongo.IndexModel{
		Keys:    bson.M{this.Sort.Field: this.Sort.Order},
		Options: options.Index().SetName(this.Name).SetUnique(this.Unique),
	}); err != nil {
		return fmt.Errorf("index complete: %w: %v(%v)", err, this.Table, this.Sort.Field)
	} // if

	return nil
}
