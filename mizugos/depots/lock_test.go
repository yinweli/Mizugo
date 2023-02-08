package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/testdata"
)

func TestLock(t *testing.T) {
	suite.Run(t, new(SuiteLock))
}

type SuiteLock struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestRedis
	major *Major
}

func (this *SuiteLock) SetupSuite() {
	this.Change("test-depots-lock")
	this.major, _ = newMajor(contexts.Ctx(), "redisdb://127.0.0.1:6379/")
}

func (this *SuiteLock) TearDownSuite() {
	this.Restore()
	this.RestoreRedis(contexts.Ctx(), this.major.Client())
	this.major.stop()
}

func (this *SuiteLock) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteLock) TestLock() {
	lock := &Lock{
		Key:  this.Key("lock"),
		time: testdata.RedisTimeout,
	}
	unlock := &Unlock{
		Key: lock.Key,
	}
	runner := this.major.Runner()

	assert.Nil(this.T(), lock.Prepare(contexts.Ctx(), runner, nil))
	_, _ = runner.Exec(contexts.Ctx())
	assert.Nil(this.T(), lock.Result())
	assert.Nil(this.T(), unlock.Prepare(contexts.Ctx(), runner, nil))
	_, _ = runner.Exec(contexts.Ctx())
	assert.Nil(this.T(), unlock.Result())

	assert.Nil(this.T(), lock.Prepare(contexts.Ctx(), runner, nil))
	_, _ = runner.Exec(contexts.Ctx())
	assert.Nil(this.T(), lock.Result())
	assert.Nil(this.T(), lock.Prepare(contexts.Ctx(), runner, nil))
	_, _ = runner.Exec(contexts.Ctx())
	assert.NotNil(this.T(), lock.Result())
}
