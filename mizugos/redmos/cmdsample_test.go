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

func TestCmdSample(t *testing.T) {
	suite.Run(t, new(SuiteCmdSample))
}

type SuiteCmdSample struct {
	suite.Suite
	trials.Catalog
	meta  testMetaSample
	major *Major
	minor *Minor
}

func (this *SuiteCmdSample) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdsample"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdsample")
}

func (this *SuiteCmdSample) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdSample) TestSample() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataSample{K: "key1", D: helps.RandStringDefault()}
	data2 := &testDataSample{K: "key2", D: helps.RandStringDefault()}

	set := &Set[testDataSample]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[testDataSample]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataSample](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey("key1"), data1))
	this.True(trials.MongoEqual[testDataSample](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey("key2"), data2))

	target := &Sample[testDataSample]{Meta: &this.meta, Count: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Len(target.Data, 1)

	target = &Sample[testDataSample]{Meta: &this.meta, Count: 2}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Len(target.Data, 2)

	target = &Sample[testDataSample]{Meta: &this.meta, Count: 5}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.LessOrEqual(len(target.Data), 5)

	target = &Sample[testDataSample]{Meta: nil, Count: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Sample[testDataSample]{Meta: &this.meta, Count: 0}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &Sample[testDataSample]{Meta: &this.meta, Count: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaSample struct {
	table bool
}

func (this *testMetaSample) MajorKey(key any) string {
	return fmt.Sprintf("cmdsample:%v", key)
}

func (this *testMetaSample) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaSample) MinorTable() string {
	if this.table {
		return "cmdsample"
	} // if

	return ""
}

type testDataSample struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
