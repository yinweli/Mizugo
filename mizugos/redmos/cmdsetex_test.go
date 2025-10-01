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

func TestCmdSetEx(t *testing.T) {
	suite.Run(t, new(SuiteCmdSetEx))
}

type SuiteCmdSetEx struct {
	suite.Suite
	trials.Catalog
	meta  testMetaSetEx
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
	data := &testDataSetEx{K: "redis+mongo", D: helps.RandStringDefault()}
	dataRedis := &testDataSetEx{K: "redis", D: helps.RandStringDefault()}
	dataMongo := &testDataSetEx{K: "mongo", D: helps.RandStringDefault()}

	target := &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target = &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: dataRedis.K, Data: dataRedis}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetEx](this.major.Client(), this.meta.MajorKey(dataRedis.K), dataRedis))
	this.False(trials.MongoEqual[testDataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataRedis.K), dataRedis))

	target = &SetEx[testDataSetEx]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: dataMongo.K, Data: dataMongo}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.False(trials.RedisEqual[testDataSetEx](this.major.Client(), this.meta.MajorKey(dataMongo.K), dataMongo))
	this.True(trials.MongoEqual[testDataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(dataMongo.K), dataMongo))

	target = &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: "", Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = true
	target = &SetEx[testDataSetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Expire: trials.Timeout * 2, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	this.True(trials.RedisEqual[testDataSetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
	trials.WaitTimeout()
	this.False(trials.RedisEqual[testDataSetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataSetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
}

func (this *SuiteCmdSetEx) TestSetExSave() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &testDataSetExSave{Save: NewSave(), K: helps.RandStringDefault(), D: helps.RandStringDefault()}
	opt := cmpopts.IgnoreFields(testDataSetExSave{}, "Save")

	data.save = true
	target := &SetEx[testDataSetExSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataSetExSave](this.major.Client(), this.meta.MajorKey(data.K), data, opt))
	this.True(trials.MongoEqual[testDataSetExSave](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data, opt))

	data.save = false
	target = &SetEx[testDataSetExSave]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
}

type testMetaSetEx struct {
	table bool
}

func (this *testMetaSetEx) MajorKey(key any) string {
	return fmt.Sprintf("cmdsetex:%v", key)
}

func (this *testMetaSetEx) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaSetEx) MinorTable() string {
	if this.table {
		return "cmdsetex"
	} // if

	return ""
}

type testDataSetEx struct {
	K string `bson:"k"`
	D string `bson:"d"`
}

type testDataSetExSave struct {
	*Save
	K string `bson:"k"`
	D string `bson:"d"`
}
