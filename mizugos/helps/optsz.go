package helps

// Optsz 選項字串
type Optsz int

// Optszf 變更選項函式類型
type Optszf func(option string) string

// On 開啟選項
func (this Optsz) On(option string) string {
	return FlagszSet(option, int32(this), true)
}

// Off 關閉選項
func (this Optsz) Off(option string) string {
	return FlagszSet(option, int32(this), false)
}

// Get 取得選項是否開啟
func (this Optsz) Get(option string) bool {
	return FlagszGet(option, int32(this))
}
