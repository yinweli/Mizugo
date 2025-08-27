package helps

import (
	"fmt"
	"slices"
)

// NewDice 建立骰子資料
func NewDice() *Dice {
	return &Dice{}
}

// Dice 骰子資料
//
// 它允許將多個元素依照「權重」放入, 之後可以隨機抽取:
//   - Rand: 依照權重隨機回傳一個元素, 不會移除
//   - RandOnce: 依照權重隨機回傳一個元素, 並從池中移除
//   - Complete: 將最後一個元素的權重補到指定總和
//   - Valid: 判斷是否還有可用的元素
//   - Max: 取得目前權重總和
//
// 基本用法:
//
//	// 建立骰子
//	d := helps.NewDice()
//
//	// 填入元素與權重
//	_ = d.One("apple", 10)  // apple 權重 10
//	_ = d.One("banana", 20) // banana 權重 20
//	_ = d.One("orange", 30) // orange 權重 30
//
//	// 擲骰 (不移除)
//	result := d.Rand()
//	fmt.Println("選到:", result)
//
//	// 擲骰 (移除)
//	for d.Valid() {
//	    fmt.Println("一次性選到:", d.RandOnce())
//	} // for
//
// 進階用法:
//
//	// 一次填入多個元素
//	_ = d.Fill([]any{"a", "b", "c"}, []int64{5, 10, 15})
//
//	// 設定總權重上限為 100, 自動補上缺少的權重
//	_ = d.Complete("other", 100)
//
// 注意:
//   - 權重必須為非負數 (weight < 0 會回錯誤)
//   - 權重 = 0 的元素會被忽略, 不會出現在結果中
//   - RandInt64n(0, this.maximum) 須保證回傳 [0, maximum) 範圍
type Dice struct {
	dice    []dice // 骰子元素列表
	maximum int64  // 最大值
}

// dice 骰子元素資料
type dice struct {
	payload any   // 酬載
	offset  int64 // 位移
}

// Clear 清除資料
func (this *Dice) Clear() {
	this.dice = []dice{}
	this.maximum = 0
}

// One 填入一筆資料
func (this *Dice) One(payload any, weight int64) error {
	if weight < 0 {
		return fmt.Errorf("dice one: weight < 0")
	} // if

	if weight == 0 {
		return nil
	} // if

	this.maximum += weight
	this.dice = append(this.dice, dice{
		payload: payload,
		offset:  this.maximum,
	})
	return nil
}

// Fill 填入多筆資料
func (this *Dice) Fill(payload []any, weight []int64) error {
	if len(payload) != len(weight) {
		return fmt.Errorf("dice fill: len mismatch")
	} // if

	for i := range payload {
		if err := this.One(payload[i], weight[i]); err != nil {
			return fmt.Errorf("dice fill: %w", err)
		} // if
	} // for

	return nil
}

// Complete 填入最後資料
func (this *Dice) Complete(payload any, maximum int64) error {
	weight := maximum - this.maximum

	if weight <= 0 {
		return nil
	} // if

	if err := this.One(payload, weight); err != nil {
		return fmt.Errorf("dice complete: %w", err)
	} // if

	return nil
}

// Rand 擲骰
func (this *Dice) Rand() any {
	if this.maximum <= 0 {
		return nil
	} // if

	num := RandInt64n(0, this.maximum)

	if find, _ := slices.BinarySearchFunc(this.dice, num, func(d dice, t int64) int {
		if d.offset >= t {
			return 1
		} // if

		return -1
	}); find < len(this.dice) {
		return this.dice[find].payload
	} // if

	return nil // 其實不管怎麼執行都不會跑到這邊
}

// RandOnce 擲骰並移除, 擲出的資料會從骰子中移除
func (this *Dice) RandOnce() any {
	if this.maximum <= 0 {
		return nil
	} // if

	num := RandInt64n(0, this.maximum)

	if find, _ := slices.BinarySearchFunc(this.dice, num, func(d dice, t int64) int {
		if d.offset >= t {
			return 1
		} // if

		return -1
	}); find < len(this.dice) {
		found := this.dice[find]
		offset := found.offset

		if find > 0 {
			offset -= this.dice[find-1].offset
		} // if

		this.dice = append(this.dice[:find], this.dice[find+1:]...)

		for n := find; n < len(this.dice); n++ {
			this.dice[n].offset -= offset
		} // for

		this.maximum -= offset
		return found.payload
	} // if

	return nil // 其實不管怎麼執行都不會跑到這邊
}

// Valid 取得骰子是否有效, 如果骰子內沒有元素, 或是所有元素的權重都是0, 則此骰子無效
func (this *Dice) Valid() bool {
	return this.maximum > 0
}

// Max 取得最大值
func (this *Dice) Max() int64 {
	return this.maximum
}
