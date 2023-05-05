package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMixed(t *testing.T) {
	suite.Run(t, new(SuiteMixed))
}

type SuiteMixed struct {
	suite.Suite
	testdata.Env
	major *Major
	minor *Minor
}

func (this *SuiteMixed) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mixed")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mixed")
}

func (this *SuiteMixed) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMixed) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMixed) TestNewMixed() {
	assert.NotNil(this.T(), newMixed(this.major, this.minor))
}

func (this *SuiteMixed) TestSubmit() {
	target := newMixed(this.major, this.minor)
	assert.NotNil(this.T(), target.Submit(ctxs.Get().Ctx()))
}

func (this *SuiteMixed) TestExec() {
	key := "mixed exec"
	target := newMixed(this.major, this.minor)
	assert.Nil(this.T(), target.Submit(ctxs.Get().Ctx()).Add(newBehaveTester(true, true)).Exec())
	assert.Nil(this.T(), target.Submit(ctxs.Get().Ctx()).Add(newBehaveTester(true, true), newBehaveTester(true, true)).Exec())
	assert.Nil(this.T(), target.Submit(ctxs.Get().Ctx()).Lock(key).Unlock(key).Exec())
	assert.Nil(this.T(), target.Submit(ctxs.Get().Ctx()).LockIf(key, true).UnlockIf(key, true).Exec())
	assert.Nil(this.T(), target.Submit(ctxs.Get().Ctx()).LockIf(key, false).UnlockIf(key, false).Exec())
	assert.NotNil(this.T(), target.Submit(ctxs.Get().Ctx()).Add(newBehaveTester(false, true)).Exec())
	assert.NotNil(this.T(), target.Submit(ctxs.Get().Ctx()).Add(newBehaveTester(true, false)).Exec())
}

func (this *SuiteMixed) TestBehave() {
	target := Behave{
		context: ctxs.Get().Ctx(),
		major:   this.major.Submit(),
		minor:   this.minor.Submit(),
	}
	assert.NotNil(this.T(), target.Ctx())
	assert.NotNil(this.T(), target.Major())
	assert.NotNil(this.T(), target.Minor())
}
