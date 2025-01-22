package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdQPush(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPush))
}

type SuiteCmdQPush struct {
	suite.Suite
	trials.Catalog
	meta  metaQPush
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
	data1 := &dataQPush{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &dataQPush{K: "redis+mongo", D: helps.RandStringDefault()}

	target := &QPush[dataQPush]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	target = &QPush[dataQPush]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompareList[dataQPush](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPush{data1, data2}))
	assert.True(this.T(), trials.MongoCompare[QueueData[dataQPush]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[dataQPush]{
		Data: []*dataQPush{data1, data2},
	}))

	target = &QPush[dataQPush]{MinorEnable: true, Meta: nil, Key: data1.K, Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &QPush[dataQPush]{MinorEnable: true, Meta: &this.meta, Key: "", Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &QPush[dataQPush]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: nil}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &QPush[dataQPush]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaQPush struct {
	table bool
}

func (this *metaQPush) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpush:%v", key)
}

func (this *metaQPush) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaQPush) MinorTable() string {
	if this.table {
		return "cmdqpush"
	} // if

	return ""
}

type dataQPush struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
