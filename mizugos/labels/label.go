package labels

import (
	"github.com/emirpasic/gods/sets/hashset"
)

// NewLabel 建立標籤資料
func NewLabel() *Label {
	return &Label{
		data: hashset.New(),
	}
}

// Label 標籤資料
type Label struct {
	data *hashset.Set // 標籤列表
}

// add 新增標籤
func (this *Label) add(label ...string) {
	for _, itor := range label {
		this.data.Add(itor)
	} // for
}

// del 刪除標籤
func (this *Label) del(label ...string) {
	for _, itor := range label {
		this.data.Remove(itor)
	} // for
}

// erase 清除標籤
func (this *Label) erase() {
	this.data = hashset.New()
}

// Label 取得標籤
func (this *Label) label() []string {
	result := []string{}

	for _, label := range this.data.Values() {
		result = append(result, label.(string))
	} // for

	return result
}
