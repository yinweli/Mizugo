package helps

import (
	"fmt"
	"sort"
)

// NewDice 建立骰子資料
func NewDice() *Dice {
	return &Dice{}
}

// Dice 骰子資料
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

// Rand 擲骰, 最大值用內部設置
func (this *Dice) Rand() any {
	return this.Randn(this.maximum)
}

// Randn 擲骰, 最大值用外部設置
func (this *Dice) Randn(maximum int64) any {
	if maximum <= 0 {
		return nil
	} // if

	num := RandInt64n(0, maximum)

	if find := sort.Search(len(this.dice), func(i int) bool {
		return this.dice[i].offset >= num
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

	if find := sort.Search(len(this.dice), func(i int) bool {
		return this.dice[i].offset >= num
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

// NewDiceDetect 建立骰子檢測資料
func NewDiceDetect() *DiceDetect {
	return &DiceDetect{
		data: map[any]int{},
	}
}

// DiceDetect 骰子檢測資料
type DiceDetect struct {
	data map[any]int // 檢測列表
}

// Add 增加檢測資料
func (this *DiceDetect) Add(key any, count int) {
	this.data[key] += count
}

// Ratio 取得檢測比例
func (this *DiceDetect) Ratio(key any, total int) float64 {
	return float64(this.data[key]) / float64(total)
}

// Check 檢測比例是否正確, 比對的方式為檢測比例是否在min與max之間
func (this *DiceDetect) Check(key any, total int, minimum, maximum float64) bool {
	ratio := this.Ratio(key, total)
	return ratio >= minimum && ratio <= maximum
}
