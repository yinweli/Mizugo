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
	major *Major
	minor *Minor
}

func (this *SuiteMxIncr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxincr")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
}

func (this *SuiteMxIncr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxIncr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxIncr) TestIncr() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()

	data := &dataMxIncr{
		Key:   "mxincr",
		Value: 2,
	}
	incr := &Incr[int64]{Meta: &this.meta, Key: data.Key, Incr: 2}
	incr.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), incr.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), incr.Complete())
	get := &Get[int64]{Meta: &this.meta, Key: data.Key}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.NotNil(this.T(), get.Data)
	assert.Equal(this.T(), data.Value, *get.Data)
	assert.True(this.T(), testdata.MongoCompare[dataMxIncr](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Key), data))

	incr = &Incr[int64]{Meta: nil, Key: data.Key, Incr: 1}
	assert.NotNil(this.T(), incr.Prepare())

	incr = &Incr[int64]{Meta: &this.meta, Key: "", Incr: 1}
	assert.NotNil(this.T(), incr.Prepare())

	this.meta.tableEmpty = false
	this.meta.fieldEmpty = true
	incr = &Incr[int64]{Meta: &this.meta, Key: data.Key, Incr: 1}
	assert.NotNil(this.T(), incr.Prepare())

	this.meta.tableEmpty = true
	this.meta.fieldEmpty = false
	incr = &Incr[int64]{Meta: &this.meta, Key: data.Key, Incr: 1}
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
