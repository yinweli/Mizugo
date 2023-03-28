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
	name  string
	major *Major
	minor *Minor
}

func (this *SuiteMixed) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-mixed")
	this.name = "mixed"
	this.major, _ = newMajor(ctxs.RootCtx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.RootCtx(), testdata.MongoURI, this.name)
}

func (this *SuiteMixed) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
	testdata.RedisClear(ctxs.RootCtx(), this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(ctxs.RootCtx(), this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.RootCtx())
}

func (this *SuiteMixed) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMixed) TestNewMixed() {
	assert.NotNil(this.T(), newMixed(this.major, this.minor))
}

func (this *SuiteMixed) TestSubmit() {
	target := newMixed(this.major, this.minor)
	assert.NotNil(this.T(), target.Submit(ctxs.RootCtx()))
}

func (this *SuiteMixed) TestExec() {
	target := newMixed(this.major, this.minor)
	key := "mixed exec"
	assert.Nil(this.T(), target.Submit(ctxs.RootCtx()).Add(newBehaveTester(true, true)).Exec())
	assert.Nil(this.T(), target.Submit(ctxs.RootCtx()).Lock(key).Unlock(key).Exec())
	assert.NotNil(this.T(), target.Submit(ctxs.RootCtx()).Add(newBehaveTester(false, true)).Exec())
	assert.NotNil(this.T(), target.Submit(ctxs.RootCtx()).Add(newBehaveTester(true, false)).Exec())
}

func (this *SuiteMixed) TestBehave() {
	target := Behave{
		context: ctxs.RootCtx(),
		major:   this.major.Submit(),
		minor:   this.minor.Submit(),
	}
	assert.NotNil(this.T(), target.Ctx())
	assert.NotNil(this.T(), target.Major())
	assert.NotNil(this.T(), target.Minor())
}
