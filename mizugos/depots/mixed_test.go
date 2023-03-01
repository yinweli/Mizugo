package depots

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
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
	name  string
	major *Major
	minor *Minor
}

func (this *SuiteMixed) SetupSuite() {
	this.Change("test-depots-mixed")
	this.name = "mixed"
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.name)
}

func (this *SuiteMixed) TearDownSuite() {
	this.Restore()
	this.RedisClear(ctxs.RootCtx(), this.major.Client())
	this.MongoClear(ctxs.RootCtx(), this.minor.Submit(this.name))
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMixed) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMixed) TestNewMixed() {
	assert.NotNil(this.T(), newMixed(this.major, this.minor))
}

func (this *SuiteMixed) TestSubmit() {
	target := newMixed(this.major, this.minor)
	assert.NotNil(this.T(), target.Submit(ctxs.Root(), this.name))
}

func (this *SuiteMixed) TestExec() {
	target := newMixed(this.major, this.minor)
	key := this.Key("mixed exec")
	assert.Nil(this.T(), target.Submit(ctxs.Root(), this.name).Add(newBehaveTester(true, true)).Exec())
	assert.Nil(this.T(), target.Submit(ctxs.Root(), this.name).Lock(key).Unlock(key).Exec())
	assert.NotNil(this.T(), target.Submit(ctxs.Root(), this.name).Add(newBehaveTester(false, true)).Exec())
	assert.NotNil(this.T(), target.Submit(ctxs.Root(), this.name).Add(newBehaveTester(true, false)).Exec())
}

func (this *SuiteMixed) TestBehave() {
	target := Behave{
		ctx:   ctxs.Root(),
		major: this.major.Submit(),
		minor: this.minor.Submit(this.name),
	}
	assert.NotNil(this.T(), target.Ctx())
	assert.NotNil(this.T(), target.Major())
	assert.NotNil(this.T(), target.Minor())
}
