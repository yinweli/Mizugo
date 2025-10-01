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

func TestCmdGetEx(t *testing.T) {
	suite.Run(t, new(SuiteCmdGetEx))
}

type SuiteCmdGetEx struct {
	suite.Suite
	trials.Catalog
	meta  testMetaGetEx
	major *Major
	minor *Minor
}

func (this *SuiteCmdGetEx) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdgetex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdgetex")
}

func (this *SuiteCmdGetEx) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdGetEx) TestGetEx() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &testDataGetEx{K: "redis+mongo", D: helps.RandStringDefault()}

	set := &Set[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisEqual[testDataGetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataGetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target := &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data, target.Data))

	target = &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data, target.Data))

	target = &GetEx[testDataGetEx]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data, target.Data))

	target = &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: testdata.Unknown}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.Nil(target.Data)

	target = &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = true
	target = &GetEx[testDataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Expire: trials.Timeout * 2, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	this.True(trials.RedisEqual[testDataGetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataGetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
	trials.WaitTimeout()
	this.False(trials.RedisEqual[testDataGetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	this.True(trials.MongoEqual[testDataGetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
}

type testMetaGetEx struct {
	table bool
}

func (this *testMetaGetEx) MajorKey(key any) string {
	return fmt.Sprintf("cmdgetex:%v", key)
}

func (this *testMetaGetEx) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaGetEx) MinorTable() string {
	if this.table {
		return "cmdgetex"
	} // if

	return ""
}

type testDataGetEx struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
