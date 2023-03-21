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
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
	dbtable string
	key     string
	major   *Major
	minor   *Minor
}

func (this *SuiteMxLock) SetupSuite() {
	this.Change("test-redmos-mxlock")
	this.dbtable = "mxlock"
	this.key = this.Key("mxlock")
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.dbtable)
}

func (this *SuiteMxLock) TearDownSuite() {
	this.Restore()
	this.RedisClear(ctxs.RootCtx(), this.major.Client())
	this.MongoClear(ctxs.RootCtx(), this.minor.Database().Collection(this.dbtable))
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMxLock) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMxLock) TestLock() {
	submit := this.major.Submit()
	lock := &Lock{Key: this.key, time: testdata.RedisTimeout}
	lock.Initialize(ctxs.Root(), submit, nil)
	unlock := &Unlock{Key: this.key}
	unlock.Initialize(ctxs.Root(), submit, nil)

	assert.Nil(this.T(), lock.Prepare())
	_, _ = submit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), lock.Complete())

	assert.Nil(this.T(), unlock.Prepare())
	_, _ = submit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), unlock.Complete())

	assert.Nil(this.T(), lock.Prepare())
	_, _ = submit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), lock.Complete())

	assert.Nil(this.T(), lock.Prepare())
	_, _ = submit.Exec(ctxs.RootCtx())
	assert.NotNil(this.T(), lock.Complete())

	lock.Key = ""
	assert.NotNil(this.T(), lock.Prepare())

	unlock.Key = ""
	assert.NotNil(this.T(), unlock.Prepare())
}
