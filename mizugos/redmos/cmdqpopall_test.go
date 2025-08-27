package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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
	meta  metaQPopAll
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
	data1 := &dataQPopAll{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &dataQPopAll{K: "redis+mongo", D: helps.RandStringDefault()}

	qpush := &QPush[dataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	qpush = &QPush[dataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompareList[dataQPopAll](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPopAll{data1, data2}))
	assert.True(this.T(), trials.MongoCompare[QueueData[dataQPopAll]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[dataQPopAll]{
		Data: []*dataQPopAll{data1, data2},
	}))

	target := &QPopAll[dataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(&QueueData[dataQPopAll]{Data: []*dataQPopAll{data1, data2}}, target.Data))
	assert.True(this.T(), trials.RedisCompareList[dataQPop](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPop{}))
	assert.True(this.T(), trials.MongoCompare[QueueData[dataQPop]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[dataQPop]{}))

	target = &QPopAll[dataQPopAll]{MinorEnable: true, Meta: nil, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &QPopAll[dataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &QPopAll[dataQPopAll]{MinorEnable: true, Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaQPopAll struct {
	table bool
}

func (this *metaQPopAll) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpopall:%v", key)
}

func (this *metaQPopAll) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaQPopAll) MinorTable() string {
	if this.table {
		return "cmdqpopall"
	} // if

	return ""
}

type dataQPopAll struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
