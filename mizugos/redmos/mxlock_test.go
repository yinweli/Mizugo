package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxLock(t *testing.T) {
	suite.Run(t, new(SuiteMxLock))
}

type SuiteMxLock struct {
	suite.Suite
	testdata.Env
	major *Major
	minor *Minor
}

func (this *SuiteMxLock) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxlock")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mxlock")
}

func (this *SuiteMxLock) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxLock) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxLock) TestLock() {
	majorSubmit := this.major.Submit()

	key := "mxlock"
	lock := &Lock{Key: key, time: testdata.RedisTimeout}
	lock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)
	unlock := &Unlock{Key: key}
	unlock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)

	assert.Nil(this.T(), lock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), lock.Complete())

	assert.Nil(this.T(), lock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.NotNil(this.T(), lock.Complete())

	assert.Nil(this.T(), unlock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), unlock.Complete())

	lock.Key = ""
	assert.NotNil(this.T(), lock.Prepare())

	unlock.Key = ""
	assert.NotNil(this.T(), unlock.Prepare())
}
