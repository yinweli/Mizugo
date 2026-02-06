package redmos

// NewSave 建立儲存判斷資料
func NewSave() *Save {
	return &Save{}
}

// Save 儲存判斷資料
//
// 用於描述當前物件是否需要儲存到主要/次要資料庫, 通常配合 Saver 介面使用, 決定物件是否需要持久化
type Save struct {
	save bool // 儲存旗標
}

// SetSave 設定儲存旗標
func (this *Save) SetSave() {
	this.save = true
}

// ClrSave 清除儲存旗標
func (this *Save) ClrSave() {
	this.save = false
}

// GetSave 取得儲存旗標
func (this *Save) GetSave() bool {
	return this.save
}
