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

func TestCmdQPeek(t *testing.T) {
	suite.Run(t, new(SuiteCmdQPeek))
}

type SuiteCmdQPeek struct {
	suite.Suite
	trials.Catalog
	meta  metaQPeek
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
	data1 := &dataQPeek{K: "redis+mongo", D: helps.RandStringDefault()}
	data2 := &dataQPeek{K: "redis+mongo", D: helps.RandStringDefault()}

	qpush := &QPush[dataQPeek]{MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	qpush = &QPush[dataQPeek]{MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	qpush.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), qpush.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), qpush.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompareList[dataQPeek](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPeek{data1, data2}))
	assert.True(this.T(), trials.MongoCompare[QueueData[dataQPeek]](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.K), &QueueData[dataQPeek]{
		Data: []*dataQPeek{data1, data2},
	}))

	target := &QPeek[dataQPeek]{Meta: &this.meta, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(&QueueData[dataQPeek]{Data: []*dataQPeek{data1, data2}}, target.Data))
	assert.True(this.T(), trials.RedisCompareList[dataQPeek](this.major.Client(), this.meta.MajorKey(data1.K), []*dataQPeek{data1, data2}))

	target = &QPeek[dataQPeek]{Meta: nil, Key: data1.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &QPeek[dataQPeek]{Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaQPeek struct {
	table bool
}

func (this *metaQPeek) MajorKey(key any) string {
	return fmt.Sprintf("cmdqpeek:%v", key)
}

func (this *metaQPeek) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaQPeek) MinorTable() string {
	if this.table {
		return "cmdqpeek"
	} // if

	return ""
}

type dataQPeek struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
