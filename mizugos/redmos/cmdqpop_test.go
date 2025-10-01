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

func TestCmdQPop(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPop))
}

type SuiteCmdQPop struct {
	suite.Suite
	trials.Catalog
	meta  testMetaQPop
	major *Major
	minor *Minor
}

func (this *SuiteCmdQPop) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdqpop"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdqpop")
}

func (this *SuiteCmdQPop) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdQPop) TestQPop() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataQPop{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &testDataQPop{K: "redis+mongo", D: helps.RandStringDefault()}

	qpush := &QPush[testDataQPop]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	qpush = &QPush[testDataQPop]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisListEqual[testDataQPop](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPop{data1, data2}))
	this.True(trials.MongoEqual[QueueData[testDataQPop]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[testDataQPop]{
		Data: []*testDataQPop{data1, data2},
	}))

	target := &QPop[testDataQPop]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal(data1, target.Data))
	this.True(trials.RedisListEqual[testDataQPop](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPop{data2}))
	this.True(trials.MongoEqual[QueueData[testDataQPop]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[testDataQPop]{
		Data: []*testDataQPop{data2},
	}))

	target = &QPop[testDataQPop]{MinorEnable: true, Meta: nil, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &QPop[testDataQPop]{MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &QPop[testDataQPop]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaQPop struct {
	table bool
}

func (this *testMetaQPop) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpop:%v", key)
}

func (this *testMetaQPop) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaQPop) MinorTable() string {
	if this.table {
		return "cmdqpop"
	} // if

	return ""
}

type testDataQPop struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
