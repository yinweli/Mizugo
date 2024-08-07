package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdSet(t *testing.T) {
	suite.Run(t, new(SuiteCmdSet))
}

type SuiteCmdSet struct {
	suite.Suite
	trials.Catalog
	meta  metaSet
	major *Major
	minor *Minor
}

func (this *SuiteCmdSet) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdset"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdset")
}

func (this *SuiteCmdSet) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSet) TestSet() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	dataAll := &dataSet{Field: "redis+mongo", Value: helps.RandStringDefault()}
	dataRedis := &dataSet{Field: "redis", Value: helps.RandStringDefault()}
	dataMongo := &dataSet{Field: "mongo", Value: helps.RandStringDefault()}

	target := &Set[dataSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSet](this.major.Client(), this.meta.MajorKey(dataAll.Field), dataAll))
	assert.True(this.T(), trials.MongoCompare[dataSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataAll.Field), dataAll))

	target = &Set[dataSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: dataRedis.Field, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSet](this.major.Client(), this.meta.MajorKey(dataRedis.Field), dataRedis))
	assert.False(this.T(), trials.MongoCompare[dataSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataRedis.Field), dataRedis))

	target = &Set[dataSet]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: dataMongo.Field, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisCompare[dataSet](this.major.Client(), this.meta.MajorKey(dataMongo.Field), dataMongo))
	assert.True(this.T(), trials.MongoCompare[dataSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataMongo.Field), dataMongo))

	target = &Set[dataSet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Set[dataSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: "", Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Set[dataSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &Set[dataSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &Set[dataSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Set[dataSet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())
}

func (this *SuiteCmdSet) TestSetSave() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetSave{Save: NewSave(), Field: helps.RandStringDefault(), Value: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(dataSetSave{}, "Save")

	data.save = true
	set := &Set[dataSetSave]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetSave](this.major.Client(), this.meta.MajorKey(data.Field), data, opt))
	assert.True(this.T(), trials.MongoCompare[dataSetSave](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, opt))

	data.save = false
	set = &Set[dataSetSave]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type metaSet struct {
	table bool
	field bool
}

func (this *metaSet) MajorKey(key any) string {
	return fmt.Sprintf("cmdset:%v", key)
}

func (this *metaSet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaSet) MinorTable() string {
	if this.table {
		return "cmdset"
	} // if

	return ""
}

func (this *metaSet) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataSet struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}

type dataSetSave struct {
	*Save
	Field string `bson:"field"`
	Value string `bson:"value"`
}
