package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestMajor(t *testing.T) {
	suite.Run(t, new(SuiteMajor))
}

type SuiteMajor struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMajor) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-major"))
}

func (this *SuiteMajor) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	assert.Nil(this.T(), target.SwitchDB(1))
	assert.NotNil(this.T(), target.SwitchDB(999999))
	target.DropDB()

	_, err = newMajor("")
	assert.NotNil(this.T(), err)

	ping, err := target.Client().Ping(context.Background()).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), "PONG", ping)

	target.stop()
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())
	assert.NotNil(this.T(), target.SwitchDB(1))
	target.DropDB()

	_, err = newMajor(testdata.RedisURIInvalid)
	assert.NotNil(this.T(), err)
}
