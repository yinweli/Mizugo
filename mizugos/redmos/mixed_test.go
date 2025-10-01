package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestMixed(t *testing.T) {
	suite.Run(t, new(SuiteMixed))
}

type SuiteMixed struct {
	suite.Suite
	trials.Catalog
	major *Major
	minor *Minor
}

func (this *SuiteMixed) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-mixed"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mixed")
}

func (this *SuiteMixed) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMixed) TestMixed() {
	target := newMixed(this.major, this.minor)
	this.NotNil(target)
	this.NotNil(target.Submit(context.Background()))
	this.NotNil(target.Major())
	this.NotNil(target.Minor())
}

func (this *SuiteMixed) TestSubmit() {
	target := newMixed(this.major, this.minor)
	this.Nil(target.Submit(context.Background()).Add(newTestBehave(true, true)).Exec())
	this.Nil(target.Submit(context.Background()).Add(newTestBehave(true, true), newTestBehave(true, true)).Exec())
	this.Nil(target.Submit(context.Background()).Lock(testdata.Unknown, testdata.Unknown).Unlock(testdata.Unknown, testdata.Unknown).Exec())
	this.Nil(target.Submit(context.Background()).LockIf(testdata.Unknown, testdata.Unknown, true).UnlockIf(testdata.Unknown, testdata.Unknown, true).Exec())
	this.Nil(target.Submit(context.Background()).LockIf(testdata.Unknown, testdata.Unknown, false).UnlockIf(testdata.Unknown, testdata.Unknown, false).Exec())
	this.NotNil(target.Submit(context.Background()).Add(newTestBehave(false, true)).Exec())
	this.NotNil(target.Submit(context.Background()).Add(newTestBehave(true, false)).Exec())
}

func (this *SuiteMixed) TestBehave() {
	target := &Behave{
		context: context.Background(),
		major:   this.major.Submit(),
		minor:   this.minor.Submit(),
	}
	this.NotNil(target.Ctx())
	this.NotNil(target.Major())
	this.NotNil(target.Minor())
}

// newTestBehave 建立測試行為
func newTestBehave(prepare, result bool) Behavior {
	return &testBehave{
		prepare: prepare,
		result:  result,
	}
}

// testBehave 測試行為
type testBehave struct {
	Behave
	prepare      bool
	result       bool
	validPrepare bool
	validResult  bool
}

func (this *testBehave) Prepare() error {
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

func (this *testBehave) Complete() error {
	if this.result == false {
		return fmt.Errorf("complete failed")
	} // if

	this.validResult = true
	return nil
}
