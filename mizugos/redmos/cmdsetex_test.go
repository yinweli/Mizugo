package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdSetEx(t *testing.T) {
	suite.Run(t, new(SuiteCmdSetEx))
}

type SuiteCmdSetEx struct {
	suite.Suite
	trials.Catalog
	meta  metaSetEx
	major *Major
	minor *Minor
}

func (this *SuiteCmdSetEx) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdsetex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdsetex")
}

func (this *SuiteCmdSetEx) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSetEx) TestSetEx() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetEx{K: "redis+mongo", D: helps.RandStringDefault()}
	dataRedis := &dataSetEx{K: "redis", D: helps.RandStringDefault()}
	dataMongo := &dataSetEx{K: "mongo", D: helps.RandStringDefault()}

	target := &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target = &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: dataRedis.K, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(dataRedis.K), dataRedis))
	assert.False(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataRedis.K), dataRedis))

	target = &SetEx[dataSetEx]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: dataMongo.K, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(dataMongo.K), dataMongo))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataMongo.K), dataMongo))

	target = &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: "", Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	target = &SetEx[dataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Expire: trials.Timeout * 2, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	assert.True(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
	trials.WaitTimeout()
	assert.False(this.T(), trials.RedisCompare[dataSetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
}

func (this *SuiteCmdSetEx) TestSetExSave() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetExSave{Save: NewSave(), K: helps.RandStringDefault(), D: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(dataSetExSave{}, "Save")

	data.save = true
	target := &SetEx[dataSetExSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetExSave](this.major.Client(), this.meta.MajorKey(data.K), data, opt))
	assert.True(this.T(), trials.MongoCompare[dataSetExSave](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data, opt))

	data.save = false
	target = &SetEx[dataSetExSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type metaSetEx struct {
	table bool
}

func (this *metaSetEx) MajorKey(key any) string {
	return fmt.Sprintf("cmdsetex:%v", key)
}

func (this *metaSetEx) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaSetEx) MinorTable() string {
	if this.table {
		return "cmdsetex"
	} // if

	return ""
}

type dataSetEx struct {
	K string `bson:"k"`
	D string `bson:"d"`
}

type dataSetExSave struct {
	*Save
	K string `bson:"k"`
	D string `bson:"d"`
}
