package labels

import (
	"github.com/emirpasic/gods/sets/hashset"
)

// 標籤資料, 用於讓物件可以被標籤管理器管理, 可以如下來使用
//   type SomeData struct {
//       Labelobj
//   }
// 或是
//   type SomeData struct {
//       *Labelobj
//   }

// NewLabelobj 建立標籤資料
func NewLabelobj() *Labelobj {
	return &Labelobj{
		data: hashset.New(),
	}
}

// Labelobj 標籤資料
type Labelobj struct {
	data *hashset.Set // 標籤列表
}

// add 新增標籤
func (this *Labelobj) add(label ...string) {
	for _, itor := range label {
		this.data.Add(itor)
	} // for
}

// del 刪除標籤
func (this *Labelobj) del(label ...string) {
	for _, itor := range label {
		this.data.Remove(itor)
	} // for
}

// erase 清除標籤
func (this *Labelobj) erase() {
	this.data = hashset.New()
}

// Label 取得標籤
func (this *Labelobj) label() []string {
	result := []string{}

	for _, label := range this.data.Values() {
		result = append(result, label.(string))
	} // for

	return result
}
