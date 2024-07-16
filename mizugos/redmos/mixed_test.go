package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Submit(context.Background()))
	assert.Equal(this.T(), this.major, target.Major())
	assert.Equal(this.T(), this.minor, target.Minor())
}

func (this *SuiteMixed) TestExec() {
	key := "mixed queue"
	target := newMixed(this.major, this.minor)
	assert.Nil(this.T(), target.Submit(context.Background()).Add(newBehaveTester(true, true)).Exec())
	assert.Nil(this.T(), target.Submit(context.Background()).Add(newBehaveTester(true, true), newBehaveTester(true, true)).Exec())
	assert.Nil(this.T(), target.Submit(context.Background()).Lock(key).Unlock(key).Exec())
	assert.Nil(this.T(), target.Submit(context.Background()).LockIf(key, true).UnlockIf(key, true).Exec())
	assert.Nil(this.T(), target.Submit(context.Background()).LockIf(key, false).UnlockIf(key, false).Exec())
	assert.NotNil(this.T(), target.Submit(context.Background()).Add(newBehaveTester(false, true)).Exec())
	assert.NotNil(this.T(), target.Submit(context.Background()).Add(newBehaveTester(true, false)).Exec())
}

func (this *SuiteMixed) TestBehave() {
	target := Behave{
		context: context.Background(),
		major:   this.major.Submit(),
		minor:   this.minor.Submit(),
	}
	assert.NotNil(this.T(), target.Ctx())
	assert.NotNil(this.T(), target.Major())
	assert.NotNil(this.T(), target.Minor())
}

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
