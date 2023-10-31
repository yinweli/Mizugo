package redmos

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxSetNX(t *testing.T) {
	suite.Run(t, new(SuiteMxSetNX))
}

type SuiteMxSetNX struct {
	suite.Suite
	testdata.Env
	meta  metaMxSetNX
	major *Major
	minor *Minor
}

func (this *SuiteMxSetNX) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxsetnx")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mxsetnx")
}

func (this *SuiteMxSetNX) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxSetNX) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxSetNX) TestSetNX() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	dataAll := &dataMxSetNX{
		Field: "mxsetnx_redis+mongo",
		Value: helps.RandStringDefault(),
	}
	dataRedis := &dataMxSetNX{
		Field: "mxsetnx_redis",
		Value: helps.RandStringDefault(),
	}
	dataMongo := &dataMxSetNX{
		Field: "mxsetnx_mongo",
		Value: helps.RandStringDefault(),
	}

	setnx := &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), setnx.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), setnx.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxSetNX](this.major.Client(), this.meta.MajorKey(dataAll.Field), dataAll, dataAll.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxSetNX](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataAll.Field), dataAll, dataAll.ignore()))

	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: dataRedis.Field, Data: dataRedis}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), setnx.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), setnx.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxSetNX](this.major.Client(), this.meta.MajorKey(dataRedis.Field), dataRedis, dataRedis.ignore()))
	assert.False(this.T(), testdata.MongoCompare[dataMxSetNX](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataRedis.Field), dataRedis, dataRedis.ignore()))

	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: dataMongo.Field, Data: dataMongo}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), setnx.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), setnx.Complete())
	assert.False(this.T(), testdata.RedisCompare[dataMxSetNX](this.major.Client(), this.meta.MajorKey(dataMongo.Field), dataMongo, dataMongo.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxSetNX](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataMongo.Field), dataMongo, dataMongo.ignore()))

	dataAll.noSave = true
	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), setnx.Prepare())
	assert.Nil(this.T(), setnx.Complete())

	dataAll.noSave = false
	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), setnx.Prepare())
	assert.NotNil(this.T(), setnx.Complete())

	dataAll.noSave = false
	setnx = &SetNX[dataMxSetNX]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), setnx.Prepare())

	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: "", Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), setnx.Prepare())

	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: nil}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), setnx.Prepare())

	this.meta.table = false
	this.meta.field = true
	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), setnx.Prepare())

	this.meta.table = true
	this.meta.field = false
	setnx = &SetNX[dataMxSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	setnx.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), setnx.Prepare())
}

type metaMxSetNX struct {
	table bool
	field bool
}

func (this *metaMxSetNX) MajorKey(key any) string {
	return fmt.Sprintf("mxsetnx:%v", key)
}

func (this *metaMxSetNX) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxSetNX) MinorTable() string {
	if this.table {
		return "mxsetnx_table"
	} // if

	return ""
}

func (this *metaMxSetNX) MinorField() string {
	if this.field {
		return "mxsetnx_field"
	} // if

	return ""
}

type dataMxSetNX struct {
	Field string `bson:"mxsetnx_field"`
	Value string `bson:"value"`

	noSave bool
}

func (this *dataMxSetNX) Save() bool {
	return this.noSave == false
}

func (this *dataMxSetNX) ignore() cmp.Option {
	return cmpopts.IgnoreFields(dataMxSetNX{}, "noSave")
}
