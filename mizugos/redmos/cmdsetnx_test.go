package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdSetNX(t *testing.T) {
	suite.Run(t, new(SuiteCmdSetNX))
}

type SuiteCmdSetNX struct {
	suite.Suite
	trials.Catalog
	meta  metaSetNX
	major *Major
	minor *Minor
}

func (this *SuiteCmdSetNX) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdsetnx"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdsetnx")
}

func (this *SuiteCmdSetNX) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSetNX) TestSetNX() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	dataAll := &dataSetNX{K: "redis+mongo", D: helps.RandStringDefault()}
	dataRedis := &dataSetNX{K: "redis", D: helps.RandStringDefault()}
	dataMongo := &dataSetNX{K: "mongo", D: helps.RandStringDefault()}

	target := &SetNX[dataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetNX](this.major.Client(), this.meta.MajorKey(dataAll.K), dataAll))
	assert.True(this.T(), trials.MongoCompare[dataSetNX](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataAll.K), dataAll))

	target = &SetNX[dataSetNX]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: dataRedis.K, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetNX](this.major.Client(), this.meta.MajorKey(dataRedis.K), dataRedis))
	assert.False(this.T(), trials.MongoCompare[dataSetNX](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataRedis.K), dataRedis))

	target = &SetNX[dataSetNX]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: dataMongo.K, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisCompare[dataSetNX](this.major.Client(), this.meta.MajorKey(dataMongo.K), dataMongo))
	assert.True(this.T(), trials.MongoCompare[dataSetNX](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataMongo.K), dataMongo))

	target = &SetNX[dataSetNX]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetNX[dataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: "", Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &SetNX[dataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &SetNX[dataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

func (this *SuiteCmdSetNX) TestSetNXSave() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataSetNXSave{Save: NewSave(), K: helps.RandStringDefault(), D: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(dataSetNXSave{}, "Save")

	data.save = true
	set := &SetNX[dataSetNXSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataSetNXSave](this.major.Client(), this.meta.MajorKey(data.K), data, opt))
	assert.True(this.T(), trials.MongoCompare[dataSetNXSave](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data, opt))

	data.save = false
	set = &SetNX[dataSetNXSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type metaSetNX struct {
	table bool
}

func (this *metaSetNX) MajorKey(key any) string {
	return fmt.Sprintf("cmdsetnx:%v", key)
}

func (this *metaSetNX) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaSetNX) MinorTable() string {
	if this.table {
		return "cmdsetnx"
	} // if

	return ""
}

type dataSetNX struct {
	K string `bson:"k"`
	D string `bson:"d"`
}

type dataSetNXSave struct {
	*Save
	K string `bson:"k"`
	D string `bson:"d"`
}
