package redmos

import (
	"fmt"
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
	testdata.Env
	meta  metaMxIncr
	key   string
	major *Major
	minor *Minor
}

func (this *SuiteMxIncr) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-mxincr")
	this.key = "mxincr-0001"
	this.major, _ = newMajor(ctxs.Root(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.Root(), testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
}

func (this *SuiteMxIncr) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
	testdata.RedisClear(ctxs.RootCtx(), this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(ctxs.RootCtx(), this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.Root())
}

func (this *SuiteMxIncr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxIncr) TestIncr() {
	expected := &dataMxIncr{
		Key:   this.key,
		Value: 2,
	}
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	get := &Get[int64]{Meta: &this.meta, Key: this.key}
	get.Initialize(ctxs.Root(), majorSubmit, minorSubmit)
	incr := &Incr[int64]{Meta: &this.meta, Key: this.key, Incr: 1}
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
	assert.Equal(this.T(), expected.Value, *get.Data)

	assert.True(this.T(), testdata.MongoCompare[dataMxIncr](ctxs.RootCtx(), this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(this.key), expected))

	incr.Meta = nil
	incr.Key = this.key
	assert.NotNil(this.T(), incr.Prepare())

	incr.Meta = &this.meta
	incr.Key = ""
	assert.NotNil(this.T(), incr.Prepare())

	incr.Meta = &this.meta
	incr.Key = this.key
	this.meta.tableEmpty = false
	this.meta.fieldEmpty = true
	assert.NotNil(this.T(), incr.Prepare())

	incr.Meta = &this.meta
	incr.Key = this.key
	this.meta.tableEmpty = true
	this.meta.fieldEmpty = false
	assert.NotNil(this.T(), incr.Prepare())
}

type metaMxIncr struct {
	tableEmpty bool
	fieldEmpty bool
}

func (this *metaMxIncr) MajorKey(key any) string {
	return fmt.Sprintf("mxincr:%v", key)
}

func (this *metaMxIncr) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxIncr) MinorTable() string {
	if this.tableEmpty == false {
		return "mxincr_table"
	} // if

	return ""
}

func (this *metaMxIncr) MinorField() string {
	if this.fieldEmpty == false {
		return "mxincr_key"
	} // if

	return ""
}

type dataMxIncr struct {
	Key   string `bson:"mxincr_key"`
	Value int64  `bson:"value"`
}
