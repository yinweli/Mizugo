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
	major *Major
	minor *Minor
}

func (this *SuiteMxGet) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mxget")
	this.major, _ = newMajor(ctxs.Get().Ctx(), testdata.RedisURI, true)
	this.minor, _ = newMinor(ctxs.Get().Ctx(), testdata.MongoURI, "mxget")
}

func (this *SuiteMxGet) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	testdata.RedisClear(this.major.Client(), this.major.UsedKey())
	testdata.MongoClear(this.minor.Database())
	this.major.stop()
	this.minor.stop(ctxs.Get().Ctx())
}

func (this *SuiteMxGet) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMxGet) TestGet() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()

	data := &testMxGet{
		Field: "mxget",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set := &Set[testMxGet]{Meta: data, Redis: false, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[testMxGet](this.major.Client(), data.MajorKey(data.Field), data, data.ignore()))
	assert.True(this.T(), testdata.MongoCompare[testMxGet](this.minor.Database(), data.MinorTable(), data.MinorField(), data.MinorKey(data.Field), data, data.ignore()))

	get := &Get[testMxGet]{Meta: data, Key: data.Field}
	get.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), get.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), get.Complete())
	assert.True(this.T(), get.Result)
	assert.True(this.T(), cmp.Equal(data, get.Data, data.ignore()))

	data = &testMxGet{
		Field: "mxget-meta-nil",
	}
	get = &Get[testMxGet]{Meta: nil, Key: data.Field}
	assert.NotNil(this.T(), get.Prepare())

	data = &testMxGet{
		Field: "mxget-key-empty",
	}
	get = &Get[testMxGet]{Meta: data, Key: ""}
	assert.NotNil(this.T(), get.Prepare())
}

func (this *SuiteMxGet) TestSet() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()

	data := &testMxGet{
		Field: "mxset-redis",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set := &Set[testMxGet]{Meta: data, Redis: true, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[testMxGet](this.major.Client(), data.MajorKey(data.Field), data, data.ignore()))
	assert.False(this.T(), testdata.MongoCompare[testMxGet](this.minor.Database(), data.MinorTable(), data.MinorField(), data.MinorKey(data.Field), data, data.ignore()))

	data = &testMxGet{
		Field: "mxset-all",
		Value: utils.RandString(testdata.RandStringLength, testdata.RandStringLetter),
	}
	set = &Set[testMxGet]{Meta: data, Redis: false, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), set.Complete())
	assert.True(this.T(), testdata.RedisCompare[testMxGet](this.major.Client(), data.MajorKey(data.Field), data, data.ignore()))
	assert.True(this.T(), testdata.MongoCompare[testMxGet](this.minor.Database(), data.MinorTable(), data.MinorField(), data.MinorKey(data.Field), data, data.ignore()))

	data = &testMxGet{
		Field:  "mxset-no-save",
		noSave: true,
	}
	set = &Set[testMxGet]{Meta: data, Redis: false, Key: data.Field, Data: data}
	set.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	assert.Nil(this.T(), set.Complete())

	data = &testMxGet{
		Field: "mxset-meta-nil",
	}
	set = &Set[testMxGet]{Meta: nil, Redis: false, Key: data.Field, Data: data}
	assert.NotNil(this.T(), set.Prepare())

	data = &testMxGet{
		Field:      "mxset-table-empty",
		tableEmpty: true,
	}
	set = &Set[testMxGet]{Meta: data, Redis: false, Key: data.Field, Data: data}
	assert.NotNil(this.T(), set.Prepare())

	data = &testMxGet{
		Field:      "mxset-field-empty",
		fieldEmpty: true,
	}
	set = &Set[testMxGet]{Meta: data, Redis: false, Key: data.Field, Data: data}
	assert.NotNil(this.T(), set.Prepare())

	data = &testMxGet{
		Field: "mxset-key-empty",
	}
	set = &Set[testMxGet]{Meta: data, Redis: false, Key: "", Data: data}
	assert.NotNil(this.T(), set.Prepare())

	data = &testMxGet{
		Field: "mxset-data-nil",
	}
	set = &Set[testMxGet]{Meta: data, Redis: false, Key: data.Field, Data: nil}
	assert.NotNil(this.T(), set.Prepare())
}

type testMxGet struct {
	Field string `bson:"mxget_field"`
	Value string `bson:"value"`

	tableEmpty bool
	fieldEmpty bool
	noSave     bool
}

func (this *testMxGet) MajorKey(key any) string {
	return fmt.Sprintf("mxget:%v", key)
}

func (this *testMxGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMxGet) MinorTable() string {
	if this.tableEmpty == false {
		return "mxget_table"
	} // if

	return ""
}

func (this *testMxGet) MinorField() string {
	if this.fieldEmpty == false {
		return "mxget_field"
	} // if

	return ""
}

func (this *testMxGet) Save() bool {
	return this.noSave == false
}

func (this *testMxGet) ignore() cmp.Option {
	return cmpopts.IgnoreFields(testMxGet{}, "tableEmpty", "fieldEmpty", "noSave")
}
