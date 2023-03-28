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
	dbtable string
	key     string
	major   *Major
	minor   *Minor
}

func (this *SuiteMxLock) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-mxlock")
	this.dbtable = "mxlock"
	this.key = "mxlock-0001"
	this.major, _ = newMajor(ctxs.RootCtx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.RootCtx(), testdata.MongoURI, this.dbtable)
}

func (this *SuiteMxLock) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
	testdata.RedisClear(ctxs.RootCtx(), this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(ctxs.RootCtx(), this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.RootCtx())
}

func (this *SuiteMxLock) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxLock) TestLock() {
	majorSubmit := this.major.Submit()
	lock := &Lock{Key: this.key, time: testdata.RedisTimeout}
	lock.Initialize(ctxs.RootCtx(), majorSubmit, nil)
	unlock := &Unlock{Key: this.key}
	unlock.Initialize(ctxs.RootCtx(), majorSubmit, nil)

	assert.Nil(this.T(), lock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), lock.Complete())

	assert.Nil(this.T(), lock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.NotNil(this.T(), lock.Complete())

	assert.Nil(this.T(), unlock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), unlock.Complete())

	lock.Key = ""
	assert.NotNil(this.T(), lock.Prepare())

	unlock.Key = ""
	assert.NotNil(this.T(), unlock.Prepare())
}
