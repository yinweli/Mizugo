package redmos

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdGet(t *testing.T) {
	suite.Run(t, new(SuiteCmdGet))
}

type SuiteCmdGet struct {
	suite.Suite
	testdata.Env
	meta  metaGet
	major *Major
	minor *Minor
}

func (this *SuiteCmdGet) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-cmdget")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdget")
}

func (this *SuiteCmdGet) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdGet) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCmdGet) TestGet() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataGet{Field: "redis+mongo", Value: helps.RandStringDefault()}

	set := &Set[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), testdata.RedisCompare[dataGet](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), testdata.MongoCompare[dataGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	target := &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: testdata.Unknown}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Data)

	target = &Get[dataGet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: ""}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Get[dataGet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())
}

type metaGet struct {
	table bool
	field bool
}

func (this *metaGet) MajorKey(key any) string {
	return fmt.Sprintf("cmdget:%v", key)
}

func (this *metaGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaGet) MinorTable() string {
	if this.table {
		return "cmdget"
	} // if

	return ""
}

func (this *metaGet) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataGet struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}
