package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdGet(t *testing.T) {
	suite.Run(t, new(SuiteCmdGet))
}

type SuiteCmdGet struct {
	suite.Suite
	trials.Catalog
	meta  testMetaGet
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
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &testDataGet{K: "redis+mongo", D: helps.RandStringDefault()}

	set := &Set[testDataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataGet](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataGet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target := &Get[testDataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data, target.Data))

	target = &Get[testDataGet]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data, target.Data))

	target = &Get[testDataGet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data, target.Data))

	target = &Get[testDataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: testdata.Unknown}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.Nil(target.Data)

	target = &Get[testDataGet]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Get[testDataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &Get[testDataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaGet struct {
	table bool
}

func (this *testMetaGet) MajorKey(key any) string {
	return fmt.Sprintf("cmdget:%v", key)
}

func (this *testMetaGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaGet) MinorTable() string {
	if this.table {
		return "cmdget"
	} // if

	return ""
}

type testDataGet struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
