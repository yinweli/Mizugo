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
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	dataAll := &dataMxSet{
		Field: "mxset_redis+mongo",
		Value: helps.RandStringDefault(),
	}
	dataRedis := &dataMxSet{
		Field: "mxset_redis",
		Value: helps.RandStringDefault(),
	}
	dataMongo := &dataMxSet{
		Field: "mxset_mongo",
		Value: helps.RandStringDefault(),
	}

	set := &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxSet](this.major.Client(), this.meta.MajorKey(dataAll.Field), dataAll, dataAll.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataAll.Field), dataAll, dataAll.ignore()))

	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: dataRedis.Field, Data: dataRedis}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxSet](this.major.Client(), this.meta.MajorKey(dataRedis.Field), dataRedis, dataRedis.ignore()))
	assert.False(this.T(), testdata.MongoCompare[dataMxSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataRedis.Field), dataRedis, dataRedis.ignore()))

	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: dataMongo.Field, Data: dataMongo}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.False(this.T(), testdata.RedisCompare[dataMxSet](this.major.Client(), this.meta.MajorKey(dataMongo.Field), dataMongo, dataMongo.ignore()))
	assert.True(this.T(), testdata.MongoCompare[dataMxSet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(dataMongo.Field), dataMongo, dataMongo.ignore()))

	dataAll.noSave = true
	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	assert.Nil(this.T(), set.Complete())

	dataAll.noSave = false
	set = &Set[dataMxSet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: "", Data: dataAll}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: nil}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	this.meta.table = false
	this.meta.field = true
	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())

	this.meta.table = true
	this.meta.field = false
	set = &Set[dataMxSet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: dataAll.Field, Data: dataAll}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), set.Prepare())
}

type metaMxSet struct {
	table bool
	field bool
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
