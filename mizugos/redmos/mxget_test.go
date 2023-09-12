package redmos

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxGet(t *testing.T) {
	suite.Run(t, new(SuiteMxGet))
}

type SuiteMxGet struct {
	suite.Suite
	testdata.Env
	meta  metaMxGet
	major *Major
	minor *Minor
}

func (this *SuiteMxGet) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxget")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mxget")
}

func (this *SuiteMxGet) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxGet) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxGet) TestGet() {
	this.meta.major = true
	this.meta.minor = true
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataMxGet{
		Field: "mxget",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set := &Set[dataMxGet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxGet](this.major.Client(), this.meta.MajorKey(data.Field), data, data.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, data.ignore()))

	get := &Get[dataMxGet]{Meta: &this.meta, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), cmp.Equal(data, get.Data, data.ignore()))

	this.meta.major = true
	this.meta.minor = false
	get = &Get[dataMxGet]{Meta: &this.meta, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), cmp.Equal(data, get.Data, data.ignore()))

	this.meta.major = false
	this.meta.minor = true
	get = &Get[dataMxGet]{Meta: &this.meta, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), cmp.Equal(data, get.Data, data.ignore()))

	this.meta.major = true
	this.meta.minor = true
	get = &Get[dataMxGet]{Meta: &this.meta, Key: testdata.Unknown}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.Nil(this.T(), get.Data)

	get = &Get[dataMxGet]{Meta: nil, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), get.Prepare())

	get = &Get[dataMxGet]{Meta: &this.meta, Key: ""}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), get.Prepare())

	this.meta.table = false
	this.meta.field = true
	get = &Get[dataMxGet]{Meta: &this.meta, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), get.Prepare())

	this.meta.table = true
	this.meta.field = false
	get = &Get[dataMxGet]{Meta: &this.meta, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), get.Prepare())
}

type metaMxGet struct {
	major bool
	minor bool
	table bool
	field bool
}

func (this *metaMxGet) Enable() (major, minor bool) {
	return this.major, this.minor
}

func (this *metaMxGet) MajorKey(key any) string {
	return fmt.Sprintf("mxget:%v", key)
}

func (this *metaMxGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxGet) MinorTable() string {
	if this.table {
		return "mxget_table"
	} // if

	return ""
}

func (this *metaMxGet) MinorField() string {
	if this.field {
		return "mxget_field"
	} // if

	return ""
}

type dataMxGet struct {
	Field string `bson:"mxget_field"`
	Value string `bson:"value"`

	noSave bool
}

func (this *dataMxGet) Save() bool {
	return this.noSave == false
}

func (this *dataMxGet) ignore() cmp.Option {
	return cmpopts.IgnoreFields(dataMxGet{}, "noSave")
}
