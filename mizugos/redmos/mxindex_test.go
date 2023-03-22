package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxIndex(t *testing.T) {
	suite.Run(t, new(SuiteMxIndex))
}

type SuiteMxIndex struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
	dbtable string
	field   string
	major   *Major
	minor   *Minor
}

func (this *SuiteMxIndex) SetupSuite() {
	this.Change("test-redmos-mxindex")
	this.dbtable = "mxindex"
	this.field = "index"
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.dbtable)
}

func (this *SuiteMxIndex) TearDownSuite() {
	this.Restore()
	this.RedisClear(ctxs.RootCtx(), this.major.Client())
	this.MongoClear(ctxs.RootCtx(), this.minor.Database().Collection(this.dbtable))
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMxIndex) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMxIndex) TestIndex() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	index := &Index{Table: this.dbtable, Field: this.field, Order: 1}
	index.Initialize(ctxs.Root(), majorSubmit, minorSubmit)

	assert.Nil(this.T(), index.Prepare())
	assert.Nil(this.T(), index.Complete())

	index.Table = ""
	index.Field = this.field
	index.Order = 1
	assert.NotNil(this.T(), index.Prepare())

	index.Table = this.dbtable
	index.Field = ""
	index.Order = 1
	assert.NotNil(this.T(), index.Prepare())

	index.Table = this.dbtable
	index.Field = this.field
	index.Order = 0
	assert.NotNil(this.T(), index.Prepare())
}
