package redmos

// Saver 儲存判斷介面
type Saver interface {
	// GetSave 取得儲存旗標
	GetSave() bool
}

// NewSave 建立儲存判斷資料
func NewSave() *Save {
	return &Save{}
}

// Save 儲存判斷資料, 用儲存旗標來判斷是否要儲存到主要/次要資料庫
type Save struct {
	save bool // 儲存旗標
}

// SetSave 設定儲存旗標
func (this *Save) SetSave() {
	this.save = true
}

// GetSave 取得儲存旗標
func (this *Save) GetSave() bool {
	return this.save
}
