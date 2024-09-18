package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdGetEx(t *testing.T) {
	suite.Run(t, new(SuiteCmdGetEx))
}

type SuiteCmdGetEx struct {
	suite.Suite
	trials.Catalog
	meta  metaGetEx
	major *Major
	minor *Minor
}

func (this *SuiteCmdGetEx) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdgetex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdgetex")
}

func (this *SuiteCmdGetEx) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdGetEx) TestGetEx() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataGetEx{Field: "redis+mongo", Value: helps.RandStringDefault()}

	set := &Set[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataGetEx](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	target := &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: testdata.Unknown}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Data)

	target = &GetEx[dataGetEx]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &GetEx[dataGetEx]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())

	this.meta.table = true
	this.meta.field = true
	target = &GetEx[dataGetEx]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Expire: trials.Timeout * 2, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	assert.True(this.T(), trials.RedisCompare[dataGetEx](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))
	trials.WaitTimeout()
	assert.False(this.T(), trials.RedisCompare[dataGetEx](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGetEx](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))
}

type metaGetEx struct {
	table bool
	field bool
}

func (this *metaGetEx) MajorKey(key any) string {
	return fmt.Sprintf("cmdgetex:%v", key)
}

func (this *metaGetEx) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaGetEx) MinorTable() string {
	if this.table {
		return "cmdgetex"
	} // if

	return ""
}

func (this *metaGetEx) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataGetEx struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}
