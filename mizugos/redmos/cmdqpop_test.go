package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdQPop(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPop))
}

type SuiteCmdQPop struct {
	suite.Suite
	trials.Catalog
	meta  metaQPop
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
	data1 := &dataQPop{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &dataQPop{K: "redis+mongo", D: helps.RandStringDefault()}

	qpush := &QPush[dataQPop]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	qpush = &QPush[dataQPop]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompareList[dataQPop](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPop{data1, data2}))
	assert.True(this.T(), trials.MongoCompare[QueueData[dataQPop]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[dataQPop]{
		Data: []*dataQPop{data1, data2},
	}))

	target := &QPop[dataQPop]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data1, target.Data))
	assert.True(this.T(), trials.RedisCompareList[dataQPop](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPop{data2}))
	assert.True(this.T(), trials.MongoCompare[QueueData[dataQPop]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[dataQPop]{
		Data: []*dataQPop{data2},
	}))

	target = &QPop[dataQPop]{MinorEnable: true, Meta: nil, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &QPop[dataQPop]{MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &QPop[dataQPop]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaQPop struct {
	table bool
}

func (this *metaQPop) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpop:%v", key)
}

func (this *metaQPop) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaQPop) MinorTable() string {
	if this.table {
		return "cmdqpop"
	} // if

	return ""
}

type dataQPop struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
