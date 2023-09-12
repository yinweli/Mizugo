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

func TestMxSet(t *testing.T) {
	suite.Run(t, new(SuiteMxSet))
}

type SuiteMxSet struct {
	suite.Suite
	testdata.Env
	meta  metaMxSet
	major *Major
	minor *Minor
}

func (this *SuiteMxSet) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxset")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mxset")
}

func (this *SuiteMxSet) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxSet) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxSet) TestSet() {
	this.meta.major = true
	this.meta.minor = true
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataMxSet{
		Field: "mxset_redis+mongo",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set := &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxSet](this.major.Client(), this.meta.MajorKey(data.Field), data, data.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, data.ignore()))

	this.meta.major = true
	this.meta.minor = false
	data = &dataMxSet{
		Field: "mxset_redis",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set = &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxSet](this.major.Client(), this.meta.MajorKey(data.Field), data, data.ignore()))
	assert.False(this.T(), testdata.MongoCompare[dataMxSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, data.ignore()))

	this.meta.major = false
	this.meta.minor = true
	data = &dataMxSet{
		Field: "mxset_mongo",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set = &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.False(this.T(), testdata.RedisCompare[dataMxSet](this.major.Client(), this.meta.MajorKey(data.Field), data, data.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, data.ignore()))

	this.meta.major = true
	this.meta.minor = true
	data.noSave = true
	set = &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	assert.Nil(this.T(), set.Complete())

	data.noSave = false
	set = &Set[dataMxSet]{Meta: nil, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	set = &Set[dataMxSet]{Meta: &this.meta, Key: "", Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	set = &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: nil}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	this.meta.table = false
	this.meta.field = true
	set = &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	this.meta.table = true
	this.meta.field = false
	set = &Set[dataMxSet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())
}

type metaMxSet struct {
	major bool
	minor bool
	table bool
	field bool
}

func (this *metaMxSet) Enable() (major, minor bool) {
	return this.major, this.minor
}

func (this *metaMxSet) MajorKey(key any) string {
	return fmt.Sprintf("mxset:%v", key)
}

func (this *metaMxSet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxSet) MinorTable() string {
	if this.table {
		return "mxset_table"
	} // if

	return ""
}

func (this *metaMxSet) MinorField() string {
	if this.field {
		return "mxset_field"
	} // if

	return ""
}

type dataMxSet struct {
	Field string `bson:"mxset_field"`
	Value string `bson:"value"`

	noSave bool
}

func (this *dataMxSet) Save() bool {
	return this.noSave == false
}

func (this *dataMxSet) ignore() cmp.Option {
	return cmpopts.IgnoreFields(dataMxSet{}, "noSave")
}
