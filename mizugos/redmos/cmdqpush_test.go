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

func TestCmdQPush(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPush))
}

type SuiteCmdQPush struct {
	suite.Suite
	trials.Catalog
	meta  testMetaQPush
	major *Major
	minor *Minor
}

func (this *SuiteCmdQPush) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdqpush"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdqpush")
}

func (this *SuiteCmdQPush) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdQPush) TestQPush() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataQPush{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &testDataQPush{K: "redis+mongo", D: helps.RandStringDefault()}

	target := &QPush[testDataQPush]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	target = &QPush[testDataQPush]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisListEqual[testDataQPush](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPush{data1, data2}))
	this.True(trials.MongoEqual[QueueData[testDataQPush]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[testDataQPush]{
		Data: []*testDataQPush{data1, data2},
	}))

	target = &QPush[testDataQPush]{MinorEnable: true, Meta: nil, Key: data1.K, Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &QPush[testDataQPush]{MinorEnable: true, Meta: &this.meta, Key: "", Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &QPush[testDataQPush]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &QPush[testDataQPush]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaQPush struct {
	table bool
}

func (this *testMetaQPush) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpush:%v", key)
}

func (this *testMetaQPush) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaQPush) MinorTable() string {
	if this.table {
		return "cmdqpush"
	} // if

	return ""
}

type testDataQPush struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
