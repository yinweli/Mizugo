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

func TestCmdQPopAll(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPopAll))
}

type SuiteCmdQPopAll struct {
	suite.Suite
	trials.Catalog
	meta  testMetaQPopAll
	major *Major
	minor *Minor
}

func (this *SuiteCmdQPopAll) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdqpopall"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdqpopall")
}

func (this *SuiteCmdQPopAll) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdQPopAll) TestQPopAll() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataQPopAll{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &testDataQPopAll{K: "redis+mongo", D: helps.RandStringDefault()}

	qpush := &QPush[testDataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	qpush = &QPush[testDataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisListEqual[testDataQPopAll](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPopAll{data1, data2}))
	this.True(trials.MongoEqual[QueueData[testDataQPopAll]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[testDataQPopAll]{
		Data: []*testDataQPopAll{data1, data2},
	}))

	target := &QPopAll[testDataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Done: func(data []*testDataQPopAll) {
		this.Equal([]*testDataQPopAll{data1, data2}, data)
	}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal([]*testDataQPopAll{data1, data2}, target.Data))
	this.True(trials.RedisListEqual[testDataQPop](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPop{}))
	this.True(trials.MongoEqual[QueueData[testDataQPop]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[testDataQPop]{}))

	target = &QPopAll[testDataQPopAll]{MinorEnable: true, Meta: nil, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &QPopAll[testDataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &QPopAll[testDataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaQPopAll struct {
	table bool
}

func (this *testMetaQPopAll) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpopall:%v", key)
}

func (this *testMetaQPopAll) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaQPopAll) MinorTable() string {
	if this.table {
		return "cmdqpopall"
	} // if

	return ""
}

type testDataQPopAll struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
