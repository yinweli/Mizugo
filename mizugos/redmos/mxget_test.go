package redmos

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
	dbtable string
	field   string
	key     string
	major   *Major
	minor   *Minor
}

type dataMxGet struct {
	Key  string `bson:"key"`
	Data string `bson:"data"`
}

func (this *SuiteMxGet) SetupSuite() {
	this.Change("test-redmos-mxget")
	this.dbtable = "mxget"
	this.field = "key"
	this.key = this.Key("mxget")
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.dbtable)
}

func (this *SuiteMxGet) TearDownSuite() {
	this.Restore()
	this.RedisClear(ctxs.RootCtx(), this.major.Client())
	this.MongoClear(ctxs.RootCtx(), this.minor.Database().Collection(this.dbtable))
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMxGet) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMxGet) TestGet() {
	expected := &dataMxGet{
		Key:  this.key,
		Data: utils.RandString(testdata.RandStringLength),
	}
	actual := &dataMxGet{}
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	get := &Get[dataMxGet]{Key: this.key}
	get.Initialize(ctxs.Root(), majorSubmit, minorSubmit)
	set := &Set[dataMxGet]{Table: this.dbtable, Field: this.field, Key: this.key, Data: expected}
	set.Initialize(ctxs.Root(), majorSubmit, minorSubmit)

	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), set.Complete())

	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), get.Result)
	assert.Equal(this.T(), set.Data, get.Data)

	assert.True(this.T(), this.MongoFindOne(ctxs.RootCtx(), this.minor.Database().Collection(this.dbtable), this.field, this.key, actual))
	assert.Equal(this.T(), expected, actual)

	get.Key = ""
	assert.NotNil(this.T(), get.Prepare())

	set.Table = ""
	set.Field = this.field
	set.Key = this.key
	set.Data = expected
	assert.NotNil(this.T(), set.Prepare())

	set.Table = this.dbtable
	set.Field = ""
	set.Key = this.key
	set.Data = expected
	assert.NotNil(this.T(), set.Prepare())

	set.Table = this.dbtable
	set.Field = this.field
	set.Key = ""
	set.Data = expected
	assert.NotNil(this.T(), set.Prepare())

	set.Table = this.dbtable
	set.Field = this.field
	set.Key = this.key
	set.Data = nil
	assert.NotNil(this.T(), set.Prepare())
}
