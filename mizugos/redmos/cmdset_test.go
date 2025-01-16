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

func TestCmdSet(t *testing.T) {
	suite.Run(t, new(SuiteCmdSet))
}

type SuiteCmdSet struct {
	suite.Suite
	trials.Catalog
	meta  metaSet
	major *Major
	minor *Minor
}

func (this *SuiteCmdSet) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdset"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdset")
}

func (this *SuiteCmdSet) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSet) TestSet() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	dataAll := &dataSet{K: "redis+mongo", D: helps.RandStringDefault()}
	dataRedis := &dataSet{K: "redis", D: helps.RandStringDefault()}
	dataMongo := &dataSet{K: "mongo", D: helps.RandStringDefault()}

	target := &Set[dataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSet](this.major.Client(), this.meta.MajorKey(dataAll.K), dataAll))
	assert.True(this.T(), trials.MongoCompare[dataSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataAll.K), dataAll))

	target = &Set[dataSet]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: dataRedis.K, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSet](this.major.Client(), this.meta.MajorKey(dataRedis.K), dataRedis))
	assert.False(this.T(), trials.MongoCompare[dataSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataRedis.K), dataRedis))

	target = &Set[dataSet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: dataMongo.K, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisCompare[dataSet](this.major.Client(), this.meta.MajorKey(dataMongo.K), dataMongo))
	assert.True(this.T(), trials.MongoCompare[dataSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataMongo.K), dataMongo))

	target = &Set[dataSet]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Set[dataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: "", Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Set[dataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &Set[dataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

func (this *SuiteCmdSet) TestSetSave() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetSave{Save: NewSave(), K: helps.RandStringDefault(), D: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(dataSetSave{}, "Save")

	data.save = true
	set := &Set[dataSetSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetSave](this.major.Client(), this.meta.MajorKey(data.K), data, opt))
	assert.True(this.T(), trials.MongoCompare[dataSetSave](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data, opt))

	data.save = false
	set = &Set[dataSetSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type metaSet struct {
	table bool
}

func (this *metaSet) MajorKey(key any) string {
	return fmt.Sprintf("cmdset:%v", key)
}

func (this *metaSet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaSet) MinorTable() string {
	if this.table {
		return "cmdset"
	} // if

	return ""
}

type dataSet struct {
	K string `bson:"k"`
	D string `bson:"d"`
}

type dataSetSave struct {
	*Save
	K string `bson:"k"`
	D string `bson:"d"`
}
