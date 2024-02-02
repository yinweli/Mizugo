package redmos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdDel(t *testing.T) {
	suite.Run(t, new(SuiteCmdDel))
}

type SuiteCmdDel struct {
	suite.Suite
	testdata.Env
	meta  metaDel
	major *Major
	minor *Minor
}

func (this *SuiteCmdDel) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-cmddel")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmddel")
}

func (this *SuiteCmdDel) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdDel) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCmdDel) TestDel() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataDel{Field: "redis+mongo", Value: helps.RandStringDefault()}

	set := &Set[dataDel]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.True(this.T(), testdata.RedisCompare[dataDel](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), testdata.MongoCompare[dataDel](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	target := &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.False(this.T(), testdata.RedisExist(this.major.Client(), this.meta.MajorKey(data.Field)))
	assert.False(this.T(), testdata.MongoExist(this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field)))

	target = &Del{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: ""}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &Del{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Del{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())
}

type metaDel struct {
	table bool
	field bool
}

func (this *metaDel) MajorKey(key any) string {
	return fmt.Sprintf("cmddel:%v", key)
}

func (this *metaDel) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaDel) MinorTable() string {
	if this.table {
		return "cmddel"
	} // if

	return ""
}

func (this *metaDel) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataDel struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}
