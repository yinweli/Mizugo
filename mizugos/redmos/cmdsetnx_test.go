package redmos

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdSetNX(t *testing.T) {
	suite.Run(t, new(SuiteCmdSetNX))
}

type SuiteCmdSetNX struct {
	suite.Suite
	trials.Catalog
	meta  metaSetNX
	major *Major
	minor *Minor
}

func (this *SuiteCmdSetNX) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdsetnx"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdsetnx")
}

func (this *SuiteCmdSetNX) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSetNX) TestSetNX() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	dataAll := &dataSetNX{Field: "redis+mongo", Value: helps.RandStringDefault()}
	dataRedis := &dataSetNX{Field: "redis", Value: helps.RandStringDefault()}
	dataMongo := &dataSetNX{Field: "mongo", Value: helps.RandStringDefault()}

	target := &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), trials.RedisCompare[dataSetNX](this.major.Client(), this.meta.MajorKey(dataAll.Field), dataAll))
	assert.True(this.T(), trials.MongoCompare[dataSetNX](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataAll.Field), dataAll))

	target = &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: dataRedis.Field, Data: dataRedis}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), trials.RedisCompare[dataSetNX](this.major.Client(), this.meta.MajorKey(dataRedis.Field), dataRedis))
	assert.False(this.T(), trials.MongoCompare[dataSetNX](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataRedis.Field), dataRedis))

	target = &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: dataMongo.Field, Data: dataMongo}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.False(this.T(), trials.RedisCompare[dataSetNX](this.major.Client(), this.meta.MajorKey(dataMongo.Field), dataMongo))
	assert.True(this.T(), trials.MongoCompare[dataSetNX](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataMongo.Field), dataMongo))

	target = &SetNX[dataSetNX]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: "", Data: dataAll}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: nil}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &SetNX[dataSetNX]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetNX[dataSetNX]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())
}

func (this *SuiteCmdSetNX) TestSetNXSave() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetNXSave{Save: NewSave(), Field: helps.RandStringDefault(), Value: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(dataSetNXSave{}, "Save")

	data.save = true
	set := &SetNX[dataSetNXSave]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), trials.RedisCompare[dataSetNXSave](this.major.Client(), this.meta.MajorKey(data.Field), data, opt))
	assert.True(this.T(), trials.MongoCompare[dataSetNXSave](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, opt))

	data.save = false
	set = &SetNX[dataSetNXSave]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
}

type metaSetNX struct {
	table bool
	field bool
}

func (this *metaSetNX) MajorKey(key any) string {
	return fmt.Sprintf("cmdsetnx:%v", key)
}

func (this *metaSetNX) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaSetNX) MinorTable() string {
	if this.table {
		return "cmdsetnx"
	} // if

	return ""
}

func (this *metaSetNX) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataSetNX struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}

type dataSetNXSave struct {
	*Save
	Field string `bson:"field"`
	Value string `bson:"value"`
}
