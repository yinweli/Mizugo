package depots

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMixed(t *testing.T) {
	suite.Run(t, new(SuiteMixed))
}

type SuiteMixed struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestRedis
	major *Major
	minor *Minor
}

func (this *SuiteMixed) SetupSuite() {
	this.Change("test-depots-mixed")
	this.major, _ = newMajor(contexts.Ctx(), "redisdb://127.0.0.1:6379/")
	this.minor, _ = newMinor(contexts.Ctx(), "mongodb://127.0.0.1:27017/")
}

func (this *SuiteMixed) TearDownSuite() {
	this.Restore()
	this.RestoreRedis(contexts.Ctx(), this.major.Client())
	this.major.stop()
	this.minor.stop(contexts.Ctx())
}

func (this *SuiteMixed) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMixed) TestNewMixed() {
	assert.NotNil(this.T(), newMixed(this.major, this.minor))
}

func (this *SuiteMixed) TestRunner() {
	target := newMixed(this.major, this.minor)
	assert.NotNil(this.T(), target.Runner(contexts.Ctx(), "", ""))
}

func (this *SuiteMixed) TestExec() {
	target := newMixed(this.major, this.minor)
	key := this.Key("lock")
	assert.Nil(this.T(), target.Runner(contexts.Ctx(), "", "").Add(newTester(true, true)).Exec())
	assert.Nil(this.T(), target.Runner(contexts.Ctx(), "", "").Lock(key).Unlock(key).Exec())
	assert.NotNil(this.T(), target.Runner(contexts.Ctx(), "", "").Add(newTester(false, true)).Exec())
	assert.NotNil(this.T(), target.Runner(contexts.Ctx(), "", "").Add(newTester(true, false)).Exec())
}

func newTester(prepare, result bool) *tester {
	return &tester{
		prepare: prepare,
		result:  result,
	}
}

type tester struct {
	prepare      bool
	result       bool
	validPrepare bool
	validResult  bool
}

func (this *tester) Prepare(ctx context.Context, majorRunner MajorRunner, minorRunner MinorRunner) error {
	if ctx == nil {
		return fmt.Errorf("ctx nil")
	} // if

	if majorRunner == nil {
		return fmt.Errorf("majorRunner nil")
	} // if

	if minorRunner == nil {
		return fmt.Errorf("minorRunner nil")
	} // if

	if this.prepare == false {
		return fmt.Errorf("prepare failed")
	} // if

	this.validPrepare = true
	return nil
}

func (this *tester) Result() error {
	if this.result == false {
		return fmt.Errorf("result failed")
	} // if

	this.validResult = true
	return nil
}
