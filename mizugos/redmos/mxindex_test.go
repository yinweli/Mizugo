package redmos

import (
	"fmt"
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
	testdata.Env
	meta  metaMxIndex
	major *Major
	minor *Minor
}

func (this *SuiteMxIndex) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxindex")
	this.major, _ = newMajor(ctxs.Get().Ctx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.Get().Ctx(), testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
}

func (this *SuiteMxIndex) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	testdata.RedisClear(ctxs.Get().Ctx(), this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(ctxs.Get().Ctx(), this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.Get().Ctx())
}

func (this *SuiteMxIndex) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxIndex) TestIndex() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	index := &Index{Meta: &this.meta, Order: 1}
	index.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)

	assert.Nil(this.T(), index.Prepare())
	assert.Nil(this.T(), index.Complete())

	index.Meta = nil
	index.Order = 1
	assert.NotNil(this.T(), index.Prepare())

	index.Meta = &this.meta
	index.Order = 0
	assert.NotNil(this.T(), index.Prepare())

	index.Meta = &this.meta
	index.Order = 1
	this.meta.tableEmpty = false
	this.meta.fieldEmpty = true
	assert.NotNil(this.T(), index.Prepare())

	index.Meta = &this.meta
	index.Order = 1
	this.meta.tableEmpty = true
	this.meta.fieldEmpty = false
	assert.NotNil(this.T(), index.Prepare())
}

type metaMxIndex struct {
	tableEmpty bool
	fieldEmpty bool
}

func (this *metaMxIndex) MajorKey(key any) string {
	return fmt.Sprintf("mxindex:%v", key)
}

func (this *metaMxIndex) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxIndex) MinorTable() string {
	if this.tableEmpty == false {
		return "mxindex_table"
	} // if

	return ""
}

func (this *metaMxIndex) MinorField() string {
	if this.fieldEmpty == false {
		return "mxindex_key"
	} // if

	return ""
}
