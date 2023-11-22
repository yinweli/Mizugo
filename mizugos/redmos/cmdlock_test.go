package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdLock(t *testing.T) {
	suite.Run(t, new(SuiteCmdLock))
}

type SuiteCmdLock struct {
	suite.Suite
	testdata.Env
	major *Major
	minor *Minor
}

func (this *SuiteCmdLock) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-cmdlock")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdlock")
}

func (this *SuiteCmdLock) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdLock) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCmdLock) TestLock() {
	majorSubmit := this.major.Submit()

	key := "lock"
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
