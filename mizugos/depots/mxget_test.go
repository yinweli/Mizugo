package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxGet(t *testing.T) {
	suite.Run(t, new(SuiteMxGet))
}

type SuiteMxGet struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
	name  string
	field string
	major *Major
	minor *Minor
}

func (this *SuiteMxGet) SetupSuite() {
	this.Change("test-depots-mxget")
	this.name = "mxget"
	this.field = "key"
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.name)
}

func (this *SuiteMxGet) TearDownSuite() {
	this.Restore()
	this.RedisClear(ctxs.RootCtx(), this.major.Client())
	this.MongoClear(ctxs.RootCtx(), this.minor.Submit(this.name))
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMxGet) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMxGet) TestGetter() {
	data := &dataTester{
		Key:  this.Key(this.name),
		Data: utils.RandString(testdata.RandStringLength),
	}
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit(this.name)
	getter := &Getter[dataTester]{}
	getter.Initialize(ctxs.Root(), majorSubmit, minorSubmit)
	setter := &Setter[dataTester]{}
	setter.Initialize(ctxs.Root(), majorSubmit, minorSubmit)

	setter.Field = this.field
	setter.Key = data.Key
	setter.Data = data
	assert.Nil(this.T(), setter.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), setter.Complete())

	getter.Key = data.Key
	assert.Nil(this.T(), getter.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), getter.Complete())
	assert.True(this.T(), getter.Result)
	assert.Equal(this.T(), data, getter.Data)

	getter.Key = ""
	assert.NotNil(this.T(), getter.Prepare())

	setter.Field = ""
	assert.NotNil(this.T(), setter.Prepare())

	setter.Field = this.field
	setter.Key = ""
	assert.NotNil(this.T(), setter.Prepare())

	setter.Field = this.field
	setter.Key = data.Key
	setter.Data = nil
	assert.NotNil(this.T(), setter.Prepare())
}
