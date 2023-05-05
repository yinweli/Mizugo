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
	this.minor, _ = newMinor(testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
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

	get = &Get[dataMxGet]{Meta: nil, Key: data.Field}
	assert.NotNil(this.T(), get.Prepare())

	get = &Get[dataMxGet]{Meta: &this.meta, Key: ""}
	assert.NotNil(this.T(), get.Prepare())
}

func (this *SuiteMxGet) TestSet() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()

	data := &dataMxGet{
		Field: "mxget_redis",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set := &Set[dataMxGet]{Meta: &this.meta, Redis: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxGet](this.major.Client(), this.meta.MajorKey(data.Field), data, data.ignore()))
	assert.False(this.T(), testdata.MongoCompare[dataMxGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, data.ignore()))

	data = &dataMxGet{
		Field: "mxget_redis+mongo",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set = &Set[dataMxGet]{Meta: &this.meta, Redis: false, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxGet](this.major.Client(), this.meta.MajorKey(data.Field), data, data.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data, data.ignore()))

	data.noSave = true
	set = &Set[dataMxGet]{Meta: &this.meta, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	assert.Nil(this.T(), set.Complete())

	data.noSave = false
	set = &Set[dataMxGet]{Meta: nil, Key: data.Field, Data: data}
	assert.NotNil(this.T(), set.Prepare())

	this.meta.tableEmpty = true
	set = &Set[dataMxGet]{Meta: &this.meta, Key: data.Field, Data: data}
	assert.NotNil(this.T(), set.Prepare())

	this.meta.tableEmpty = false
	this.meta.fieldEmpty = true
	set = &Set[dataMxGet]{Meta: &this.meta, Key: data.Field, Data: data}
	assert.NotNil(this.T(), set.Prepare())

	this.meta.fieldEmpty = false
	set = &Set[dataMxGet]{Meta: &this.meta, Key: "", Data: data}
	assert.NotNil(this.T(), set.Prepare())

	set = &Set[dataMxGet]{Meta: &this.meta, Key: data.Field, Data: nil}
	assert.NotNil(this.T(), set.Prepare())
}

type metaMxGet struct {
	tableEmpty bool
	fieldEmpty bool
}

func (this *metaMxGet) MajorKey(key any) string {
	return fmt.Sprintf("mxget:%v", key)
}

func (this *metaMxGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxGet) MinorTable() string {
	if this.tableEmpty == false {
		return "mxget_table"
	} // if

	return ""
}

func (this *metaMxGet) MinorField() string {
	if this.fieldEmpty == false {
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
