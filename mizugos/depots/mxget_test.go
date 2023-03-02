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

func (this *SuiteMxGet) TestGet() {
	data := &dataTester{
		Key:  this.Key(this.name),
		Data: utils.RandString(testdata.RandStringLength),
	}
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit(this.name)
	get := &Get[dataTester]{}
	get.Initialize(ctxs.Root(), majorSubmit, minorSubmit)
	set := &Set[dataTester]{}
	set.Initialize(ctxs.Root(), majorSubmit, minorSubmit)

	set.Field = this.field
	set.Key = data.Key
	set.Data = data
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), set.Complete())

	get.Key = data.Key
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), get.Result)
	assert.Equal(this.T(), data, get.Data)

	get.Key = ""
	assert.NotNil(this.T(), get.Prepare())

	set.Field = ""
	set.Key = data.Key
	set.Data = data
	assert.NotNil(this.T(), set.Prepare())

	set.Field = this.field
	set.Key = ""
	set.Data = data
	assert.NotNil(this.T(), set.Prepare())

	set.Field = this.field
	set.Key = data.Key
	set.Data = nil
	assert.NotNil(this.T(), set.Prepare())
}
