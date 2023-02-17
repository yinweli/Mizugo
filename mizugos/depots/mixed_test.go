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
	testdata.TestDB
	major *Major
	minor *Minor
	name  string
}

func (this *SuiteMixed) SetupSuite() {
	this.Change("test-depots-mixed")
	this.major, _ = newMajor(contexts.Ctx(), testdata.RedisURI)
	this.minor, _ = newMinor(contexts.Ctx(), testdata.MongoURI)
	this.name = "mixed"
}

func (this *SuiteMixed) TearDownSuite() {
	this.Restore()
	this.RedisClear(contexts.Ctx(), this.major.Client())
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
	assert.NotNil(this.T(), target.Runner(contexts.Ctx(), this.name, this.name))
}

func (this *SuiteMixed) TestExec() {
	target := newMixed(this.major, this.minor)
	key := this.Key("lock")
	assert.Nil(this.T(), target.Runner(contexts.Ctx(), this.name, this.name).Add(newMixedTester(true, true)).Exec())
	assert.Nil(this.T(), target.Runner(contexts.Ctx(), this.name, this.name).Lock(key).Unlock(key).Exec())
	assert.NotNil(this.T(), target.Runner(contexts.Ctx(), this.name, this.name).Add(newMixedTester(false, true)).Exec())
	assert.NotNil(this.T(), target.Runner(contexts.Ctx(), this.name, this.name).Add(newMixedTester(true, false)).Exec())
}

func newMixedTester(prepare, result bool) *mixedTester {
	return &mixedTester{
		prepare: prepare,
		result:  result,
	}
}

type mixedTester struct {
	prepare      bool
	result       bool
	validPrepare bool
	validResult  bool
}

func (this *mixedTester) Prepare(ctx context.Context, majorRunner MajorRunner, minorRunner MinorRunner) error {
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

func (this *mixedTester) Result() error {
	if this.result == false {
		return fmt.Errorf("result failed")
	} // if

	this.validResult = true
	return nil
}
