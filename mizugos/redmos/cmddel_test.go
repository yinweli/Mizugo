package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdDel(t *testing.T) {
	suite.Run(t, new(SuiteCmdDel))
}

type SuiteCmdDel struct {
	suite.Suite
	trials.Catalog
	meta  testMetaDel
	major *Major
	minor *Minor
}

func (this *SuiteCmdDel) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmddel"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmddel")
}

func (this *SuiteCmdDel) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdDel) TestDel() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &testDataDel{K: "redis+mongo", D: helps.RandStringDefault()}

	set := &Set[testDataDel]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataDel](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataDel](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target := &Del{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.False(trials.RedisExist(this.major.Client(), this.meta.MajorKey(data.K)))
	this.False(trials.MongoExist(this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K)))

	target = &Del{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Del{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &Del{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaDel struct {
	table bool
}

func (this *testMetaDel) MajorKey(key any) string {
	return fmt.Sprintf("cmddel:%v", key)
}

func (this *testMetaDel) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaDel) MinorTable() string {
	if this.table {
		return "cmddel"
	} // if

	return ""
}

type testDataDel struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
