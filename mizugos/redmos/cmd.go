package redmos

// Saver 儲存判斷介面
type Saver interface {
	// Get 取得儲存旗標
	Get() bool
}

// Save 儲存判斷資料, 用儲存旗標來判斷是否要儲存到主要/次要資料庫
type Save struct {
	save bool // 儲存旗標
}

// Set 設定儲存旗標
func (this *Save) Set() {
	this.save = true
}

// Get 取得儲存旗標
func (this *Save) Get() bool {
	return this.save
}
