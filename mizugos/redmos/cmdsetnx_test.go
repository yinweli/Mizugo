package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
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
	meta  testMetaSetNX
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
	dataAll := &testDataSetNX{K: "redis+mongo", D: helps.RandStringDefault()}
	dataRedis := &testDataSetNX{K: "redis", D: helps.RandStringDefault()}
	dataMongo := &testDataSetNX{K: "mongo", D: helps.RandStringDefault()}

	target := &SetNX[testDataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetNX](this.major.Client(), this.meta.MajorKey(dataAll.K), dataAll))
	this.True(trials.MongoEqual[testDataSetNX](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataAll.K), dataAll))

	target = &SetNX[testDataSetNX]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: dataRedis.K, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetNX](this.major.Client(), this.meta.MajorKey(dataRedis.K), dataRedis))
	this.False(trials.MongoEqual[testDataSetNX](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataRedis.K), dataRedis))

	target = &SetNX[testDataSetNX]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: dataMongo.K, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.False(trials.RedisEqual[testDataSetNX](this.major.Client(), this.meta.MajorKey(dataMongo.K), dataMongo))
	this.True(trials.MongoEqual[testDataSetNX](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataMongo.K), dataMongo))

	target = &SetNX[testDataSetNX]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &SetNX[testDataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: "", Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &SetNX[testDataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &SetNX[testDataSetNX]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

func (this *SuiteCmdSetNX) TestSetNXSave() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &testDataSetNXSave{Save: NewSave(), K: helps.RandStringDefault(), D: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(testDataSetNXSave{}, "Save")

	data.save = true
	set := &SetNX[testDataSetNXSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetNXSave](this.major.Client(), this.meta.MajorKey(data.K), data, opt))
	this.True(trials.MongoEqual[testDataSetNXSave](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data, opt))

	data.save = false
	set = &SetNX[testDataSetNXSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type testMetaSetNX struct {
	table bool
}

func (this *testMetaSetNX) MajorKey(key any) string {
	return fmt.Sprintf("cmdsetnx:%v", key)
}

func (this *testMetaSetNX) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaSetNX) MinorTable() string {
	if this.table {
		return "cmdsetnx"
	} // if

	return ""
}

type testDataSetNX struct {
	K string `bson:"k"`
	D string `bson:"d"`
}

type testDataSetNXSave struct {
	*Save
	K string `bson:"k"`
	D string `bson:"d"`
}
