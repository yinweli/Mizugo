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
	key   string
	major *Major
	minor *Minor
}

func (this *SuiteMxGet) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-mxget")
	this.key = "mxget-0001"
	this.major, _ = newMajor(ctxs.Get().Ctx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.Get().Ctx(), testdata.MongoURI, this.meta.MinorTable()) // 這裡偷懶把表格名稱當資料庫名稱來用
}

func (this *SuiteMxGet) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
	testdata.RedisClear(ctxs.Get().Ctx(), this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(ctxs.Get().Ctx(), this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.Get().Ctx())
}

func (this *SuiteMxGet) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxGet) TestGet() {
	expected := &dataMxGet{
		dirty: true,
		Key:   this.key,
		Data:  utils.RandString(testdata.RandStringLength),
	}
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	get := &Get[dataMxGet]{Meta: &this.meta, Key: this.key}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	set := &Set[dataMxGet]{Meta: &this.meta, Key: this.key, Data: expected}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	ignore := cmpopts.IgnoreFields(dataMxGet{}, "dirty") // 比對時忽略dirty欄位

	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())

	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), get.Result)
	assert.True(this.T(), cmp.Equal(set.Data, get.Data, ignore))

	assert.True(this.T(), testdata.RedisCompare[dataMxGet](ctxs.Get().Ctx(), this.major.Client(), this.meta.MajorKey(this.key), expected, ignore))
	assert.True(this.T(), testdata.MongoCompare[dataMxGet](ctxs.Get().Ctx(), this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(this.key), expected, ignore))

	set.Data.dirty = false
	assert.Nil(this.T(), set.Prepare())
	assert.Nil(this.T(), set.Complete())

	get.Meta = nil
	get.Key = this.key
	assert.NotNil(this.T(), get.Prepare())

	get.Meta = &this.meta
	get.Key = ""
	assert.NotNil(this.T(), get.Prepare())

	get.Meta = &this.meta
	get.Key = this.key
	this.meta.tableEmpty = false
	this.meta.fieldEmpty = true
	assert.NotNil(this.T(), get.Prepare())

	get.Meta = &this.meta
	get.Key = this.key
	this.meta.tableEmpty = true
	this.meta.fieldEmpty = false
	assert.NotNil(this.T(), get.Prepare())

	set.Meta = nil
	set.Key = this.key
	set.Data = expected
	this.meta.tableEmpty = false
	this.meta.fieldEmpty = false
	assert.NotNil(this.T(), set.Prepare())

	set.Meta = &this.meta
	set.Key = ""
	set.Data = expected
	assert.NotNil(this.T(), set.Prepare())

	set.Meta = &this.meta
	set.Key = this.key
	set.Data = nil
	assert.NotNil(this.T(), set.Prepare())

	set.Meta = &this.meta
	set.Key = this.key
	set.Data = expected
	this.meta.tableEmpty = false
	this.meta.fieldEmpty = true
	assert.NotNil(this.T(), set.Prepare())

	set.Meta = &this.meta
	set.Key = this.key
	set.Data = expected
	this.meta.tableEmpty = true
	this.meta.fieldEmpty = false
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
		return "mxget_key"
	} // if

	return ""
}

type dataMxGet struct {
	dirty bool
	Key   string `bson:"mxget_key"`
	Data  string `bson:"data"`
}

func (this *dataMxGet) Save() bool {
	return this.dirty
}
