package helps

// Optsz 選項字串, 用來操作 option string (旗標字串)
//
// 旗標字串可視為一條位元列, 例如 "1010..."
// 每個 Optsz 對應到某個索引位置:
//   - On: 將該索引位置設為 1 (開啟)
//   - Off: 將該索引位置設為 0 (關閉)
//   - Get: 讀取該索引位置是否為 1 (已開啟)
//
// 範例:
//
//	const (
//	    OptTutorial Optsz = iota
//	    OptGreatWin
//	)
//
//	var option string
//	option = OptTutorial.On(option)  // 開啟教學完成
//	option = OptGreatWin.Off(option) // 關閉雙倍獎勵
//
//	if OptGreatWin.Get(option) {
//	    // 雙倍獎勵已開啟
//	} // if
type Optsz int

// Optszf 變更選項函式類型
//
// 它的存在意義是: 將對 option string 的某個操作包裝成函式, 方便做「鏈式套用」或「批次套用」, 而不是一條一條硬寫 On/Off
//
// 典型用途是把 Optsz 的 On/Off 等操作包裝起來, 讓多個操作可以重複使用或組合
//
// 範例:
//
//	const (
//	    OptTutorial Optsz = iota
//	    OptGreatWin
//	)
//
//	// 建立操作器
//	onTutorial := func(opt string) string { return OptTutorial.On(opt) }
//	offGreatWin  := func(opt string) string { return OptGreatWin.Off(opt) }
//
//	// 應用操作器
//	var option string
//	option = onTutorial(option)  // 開啟教學完成
//	option = offGreatWin(option) // 關閉雙倍獎勵
//
//	// 也可以一次套用多個操作器
//	for _, itor := range []Optszf{onTutorial, offGreatWin} {
//	    option = itor(option)
//	} // for
//
// 這樣的設計讓操作可組合, 可重複使用, 並且能被其他函式(例如批次設定器)直接接受
type Optszf func(option string) string

// On 開啟選項
func (this Optsz) On(option string) string {
	return FlagszSet(option, int(this), true)
}

// Off 關閉選項
func (this Optsz) Off(option string) string {
	return FlagszSet(option, int(this), false)
}

// Get 取得選項是否開啟
func (this Optsz) Get(option string) bool {
	return FlagszGet(option, int(this))
}
