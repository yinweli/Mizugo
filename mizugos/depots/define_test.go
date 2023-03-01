package depots

import (
	"fmt"
)

// newBehaveTester 建立測試行為
func newBehaveTester(prepare, result bool) Behavior {
	return &behaveTester{
		prepare: prepare,
		result:  result,
	}
}

// behaveTester 測試行為
type behaveTester struct {
	Behave
	prepare      bool
	result       bool
	validPrepare bool
	validResult  bool
}

func (this *behaveTester) Prepare() error {
	if this.Ctx() == nil {
		return fmt.Errorf("ctx nil")
	} // if

	if this.Major() == nil {
		return fmt.Errorf("major nil")
	} // if

	if this.Minor() == nil {
		return fmt.Errorf("minor nil")
	} // if

	if this.prepare == false {
		return fmt.Errorf("prepare failed")
	} // if

	this.validPrepare = true
	return nil
}

func (this *behaveTester) Complete() error {
	if this.result == false {
		return fmt.Errorf("complete failed")
	} // if

	this.validResult = true
	return nil
}

// dataTester 測試資料
type dataTester struct {
	Key  string `bson:"key"`
	Data string `bson:"data"`
}
