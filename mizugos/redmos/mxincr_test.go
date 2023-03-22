package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxIncr(t *testing.T) {
	suite.Run(t, new(SuiteMxIncr))
}

type SuiteMxIncr struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestDB
	dbtable string
	field   string
	key     string
	major   *Major
	minor   *Minor
}

type dataMxIncr struct {
	Name  string `bson:"name"`
	Value int64  `bson:"value"`
}

func (this *SuiteMxIncr) SetupSuite() {
	this.TBegin("test-redmos-mxincr", "")
	this.dbtable = "mxincr"
	this.field = "name"
	this.key = this.Key("mxincr")
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.dbtable)
}

func (this *SuiteMxIncr) TearDownSuite() {
	this.TFinal()
	this.RedisClear(ctxs.RootCtx(), this.major.Client())
	this.MongoClear(ctxs.RootCtx(), this.minor.Database().Collection(this.dbtable))
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMxIncr) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteMxIncr) TestIncr() {
	expected := &dataMxIncr{
		Name:  this.key,
		Value: 2,
	}
	actual := &dataMxIncr{}
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	get := &Get[int64]{Key: this.key}
	get.Initialize(ctxs.Root(), majorSubmit, minorSubmit)
	incr := &Incr[int64]{Table: this.dbtable, Field: this.field, Key: this.key, Incr: 1}
	incr.Initialize(ctxs.Root(), majorSubmit, minorSubmit)

	assert.Nil(this.T(), incr.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), incr.Complete())

	assert.Nil(this.T(), incr.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), incr.Complete())

	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.RootCtx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), get.Result)
	assert.NotNil(this.T(), get.Data)
	assert.Equal(this.T(), int64(2), *get.Data)

	assert.True(this.T(), this.MongoFindOne(ctxs.RootCtx(), this.minor.Database().Collection(this.dbtable), this.field, this.key, actual))
	assert.Equal(this.T(), expected, actual)

	incr.Table = ""
	incr.Field = this.field
	incr.Key = this.key
	assert.NotNil(this.T(), incr.Prepare())

	incr.Table = this.dbtable
	incr.Field = ""
	incr.Key = this.key
	assert.NotNil(this.T(), incr.Prepare())

	incr.Table = this.dbtable
	incr.Field = this.field
	incr.Key = ""
	assert.NotNil(this.T(), incr.Prepare())
}
