package redmos

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdGet(t *testing.T) {
	suite.Run(t, new(SuiteCmdGet))
}

type SuiteCmdGet struct {
	suite.Suite
	trials.Catalog
	meta  metaGet
	major *Major
	minor *Minor
}

func (this *SuiteCmdGet) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdget"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdget")
}

func (this *SuiteCmdGet) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdGet) TestGet() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataGet{Field: "redis+mongo", Value: helps.RandStringDefault()}

	set := &Set[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataGet](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	target := &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: false, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: testdata.Unknown}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Data)

	target = &Get[dataGet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Get[dataGet]{Meta: nil, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())

	this.meta.table = true
	this.meta.field = true
	target = &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Expire: trials.Timeout * 2}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	assert.True(this.T(), trials.RedisCompare[dataGet](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))
	trials.WaitTimeout()
	assert.False(this.T(), trials.RedisCompare[dataGet](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))
}

func (this *SuiteCmdGet) TestGetRefreshExpireTime() {
	// 過期時間 2 秒
	expiredTime := 2 * time.Second

	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataGet{Field: "redis+mongo+expire", Value: helps.RandStringDefault()}

	set := &Set[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: true, Key: data.Field, Data: data, Expire: expiredTime}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataGet](this.major.Client(), this.meta.MajorKey(data.Field), data))
	assert.True(this.T(), trials.MongoCompare[dataGet](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	// 可以撐過 3 秒代表 get 的時候有刷新 expire time
	for i := 0; i < 3; i++ {
		// 停止 1 秒
		time.Sleep(time.Second)

		// 讀取 redis 並刷新 expire time
		target := &Get[dataGet]{Meta: &this.meta, MajorEnable: true, MinorEnable: false, Key: data.Field, Expire: expiredTime}
		target.Initialize(context.Background(), majorSubmit, minorSubmit)
		assert.Nil(this.T(), target.Prepare())
		_, _ = majorSubmit.Exec(context.Background())
		assert.Nil(this.T(), target.Complete())
		_ = minorSubmit.Exec(context.Background())
		assert.True(this.T(), cmp.Equal(data, target.Data))
	}
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
