package redmos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMxDel(t *testing.T) {
	suite.Run(t, new(SuiteMxDel))
}

type SuiteMxDel struct {
	suite.Suite
	testdata.Env
	meta  metaMxDel
	major *Major
	minor *Minor
}

func (this *SuiteMxDel) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxdel")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "mxdel")
}

func (this *SuiteMxDel) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteMxDel) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxDel) TestDel() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataMxDel{
		Field: "mxdel_redis+mongo",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}

	set := &Set[dataMxDel]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[dataMxDel](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), testdata.MongoCompare[dataMxDel](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	del := &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	del.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), del.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), del.Complete())
	assert.False(this.T(), testdata.RedisExist(this.major.Client(), this.meta.MajorKey(data.Field)))
	assert.False(this.T(), testdata.MongoExist(this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field)))

	del = &Del{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field}
	del.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), del.Prepare())

	del = &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: ""}
	del.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), del.Prepare())

	this.meta.table = false
	this.meta.field = true
	del = &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	del.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), del.Prepare())

	this.meta.table = true
	this.meta.field = false
	del = &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	del.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), del.Prepare())
}

type metaMxDel struct {
	table bool
	field bool
}

func (this *metaMxDel) MajorKey(key any) string {
	return fmt.Sprintf("mxdel:%v", key)
}

func (this *metaMxDel) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaMxDel) MinorTable() string {
	if this.table {
		return "mxdel_table"
	} // if

	return ""
}

func (this *metaMxDel) MinorField() string {
	if this.field {
		return "mxdel_field"
	} // if

	return ""
}

type dataMxDel struct {
	Field string `bson:"mxdel_field"`
	Value string `bson:"value"`
}
