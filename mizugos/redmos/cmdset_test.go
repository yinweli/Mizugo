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

func TestCmdSet(t *testing.T) {
	suite.Run(t, new(SuiteCmdSet))
}

type SuiteCmdSet struct {
	suite.Suite
	trials.Catalog
	meta  testMetaSet
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
	dataAll := &testDataSet{K: "redis+mongo", D: helps.RandStringDefault()}
	dataRedis := &testDataSet{K: "redis", D: helps.RandStringDefault()}
	dataMongo := &testDataSet{K: "mongo", D: helps.RandStringDefault()}

	target := &Set[testDataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSet](this.major.Client(), this.meta.MajorKey(dataAll.K), dataAll))
	this.True(trials.MongoEqual[testDataSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataAll.K), dataAll))

	target = &Set[testDataSet]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: dataRedis.K, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSet](this.major.Client(), this.meta.MajorKey(dataRedis.K), dataRedis))
	this.False(trials.MongoEqual[testDataSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataRedis.K), dataRedis))

	target = &Set[testDataSet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: dataMongo.K, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.False(trials.RedisEqual[testDataSet](this.major.Client(), this.meta.MajorKey(dataMongo.K), dataMongo))
	this.True(trials.MongoEqual[testDataSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataMongo.K), dataMongo))

	target = &Set[testDataSet]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Set[testDataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: "", Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Set[testDataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &Set[testDataSet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: dataAll.K, Data: dataAll}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

func (this *SuiteCmdSet) TestSetSave() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &testDataSetSave{Save: NewSave(), K: helps.RandStringDefault(), D: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(testDataSetSave{}, "Save")

	data.save = true
	set := &Set[testDataSetSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetSave](this.major.Client(), this.meta.MajorKey(data.K), data, opt))
	this.True(trials.MongoEqual[testDataSetSave](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data, opt))

	data.save = false
	set = &Set[testDataSetSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type testMetaSet struct {
	table bool
}

func (this *testMetaSet) MajorKey(key any) string {
	return fmt.Sprintf("cmdset:%v", key)
}

func (this *testMetaSet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaSet) MinorTable() string {
	if this.table {
		return "cmdset"
	} // if

	return ""
}

type testDataSet struct {
	K string `bson:"k"`
	D string `bson:"d"`
}

type testDataSetSave struct {
	*Save
	K string `bson:"k"`
	D string `bson:"d"`
}
