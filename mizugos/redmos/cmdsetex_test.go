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

func TestCmdSetEx(t *testing.T) {
	suite.Run(t, new(SuiteCmdSetEx))
}

type SuiteCmdSetEx struct {
	suite.Suite
	trials.Catalog
	meta  metaSetEx
	major *Major
	minor *Minor
}

func (this *SuiteCmdSetEx) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdsetex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdsetex")
}

func (this *SuiteCmdSetEx) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSetEx) TestSetEx() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetEx{Field: "redis+mongo", Value: helps.RandStringDefault()}
	dataRedis := &dataSetEx{Field: "redis", Value: helps.RandStringDefault()}
	dataMongo := &dataSetEx{Field: "mongo", Value: helps.RandStringDefault()}

	target := &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: dataRedis.Field, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(dataRedis.Field), dataRedis))
	assert.False(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataRedis.Field), dataRedis))

	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: dataMongo.Field, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(dataMongo.Field), dataMongo))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataMongo.Field), dataMongo))

	target = &SetEx[dataSetEx]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: "", Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetEx[dataSetEx]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())

	this.meta.table = true
	this.meta.field = true
	target = &SetEx[dataSetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Expire: trials.Timeout * 2, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	assert.True(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))
	trials.WaitTimeout()
	assert.False(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))
}

func (this *SuiteCmdSetEx) TestSetExSave() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetExSave{Save: NewSave(), Field: helps.RandStringDefault(), Value: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(dataSetExSave{}, "Save")

	data.save = true
	target := &SetEx[dataSetExSave]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetExSave](this.major.Client(), this.meta.MajorKey(data.Field), data, opt))
	assert.True(this.T(), trials.MongoCompare[dataSetExSave](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, opt))

	data.save = false
	target = &SetEx[dataSetExSave]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type metaSetEx struct {
	table bool
	field bool
}

func (this *metaSetEx) MajorKey(key any) string {
	return fmt.Sprintf("cmdsetex:%v", key)
}

func (this *metaSetEx) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaSetEx) MinorTable() string {
	if this.table {
		return "cmdsetex"
	} // if

	return ""
}

func (this *metaSetEx) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataSetEx struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}

type dataSetExSave struct {
	*Save
	Field string `bson:"field"`
	Value string `bson:"value"`
}
