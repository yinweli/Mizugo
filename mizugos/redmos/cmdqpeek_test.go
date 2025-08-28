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

func TestCmdQPeek(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPeek))
}

type SuiteCmdQPeek struct {
	suite.Suite
	trials.Catalog
	meta  testMetaQPeek
	major *Major
	minor *Minor
}

func (this *SuiteCmdQPeek) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdqpeek"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdqpeek")
}

func (this *SuiteCmdQPeek) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdQPeek) TestQPeek() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataQPeek{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &testDataQPeek{K: "redis+mongo", D: helps.RandStringDefault()}

	qpush := &QPush[testDataQPeek]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	qpush = &QPush[testDataQPeek]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.RedisListEqual[testDataQPeek](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPeek{data1, data2}))
	this.True(trials.MongoEqual[QueueData[testDataQPeek]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[testDataQPeek]{
		Data: []*testDataQPeek{data1, data2},
	}))

	target := &QPeek[testDataQPeek]{Meta: &this.meta, Key: data1.K, Done: func(data []*testDataQPeek) {
		this.Equal([]*testDataQPeek{data1, data2}, data)
	}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(cmp.Equal([]*testDataQPeek{data1, data2}, target.Data))
	this.True(trials.RedisListEqual[testDataQPeek](this.major.Client(), this.meta.MajorKey(data1.K), []*testDataQPeek{data1, data2}))

	target = &QPeek[testDataQPeek]{Meta: nil, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &QPeek[testDataQPeek]{Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaQPeek struct {
	table bool
}

func (this *testMetaQPeek) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpeek:%v", key)
}

func (this *testMetaQPeek) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaQPeek) MinorTable() string {
	if this.table {
		return "cmdqpeek"
	} // if

	return ""
}

type testDataQPeek struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
