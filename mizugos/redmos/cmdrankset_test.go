package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdRankSet(t *testing.T) {
	suite.Run(t, new(SuiteCmdRankSet))
}

type SuiteCmdRankSet struct {
	suite.Suite
	trials.Catalog
	meta  testMetaRankSet
	major *Major
	minor *Minor
}

func (this *SuiteCmdRankSet) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdrankset"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdrankset")
}

func (this *SuiteCmdRankSet) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdRankSet) TestRankSet() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	key := helps.RandStringDefault()
	data1 := &testDataRankSet{K: key, Score: 100}
	data2 := &testDataRankSet{K: key, Score: 50}
	data3 := &testDataRankSet{K: key, Score: 200}

	target := &RankSet[testDataRankSet]{
		Meta:   &this.meta,
		Key:    key,
		Filter: bson.M{"score": bson.M{"$lt": 100}},
		Data:   data1,
	}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataRankSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(key), data1))

	target = &RankSet[testDataRankSet]{
		Meta:   &this.meta,
		Key:    key,
		Filter: bson.M{"score": bson.M{"$lt": 50}},
		Data:   data2,
	}
	target.Initialize(context.Background(), nil, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataRankSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(key), data1))

	target = &RankSet[testDataRankSet]{
		Meta:   &this.meta,
		Key:    key,
		Filter: bson.M{"score": bson.M{"$lt": 200}},
		Data:   data3,
	}
	target.Initialize(context.Background(), nil, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataRankSet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(key), data3))

	target = &RankSet[testDataRankSet]{Meta: nil, Key: key, Data: data1}
	target.Initialize(context.Background(), nil, minorSubmit)
	this.NotNil(target.Prepare())

	target = &RankSet[testDataRankSet]{Meta: &this.meta, Key: "", Data: data1}
	target.Initialize(context.Background(), nil, minorSubmit)
	this.NotNil(target.Prepare())

	target = &RankSet[testDataRankSet]{Meta: &this.meta, Key: key, Data: nil}
	target.Initialize(context.Background(), nil, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &RankSet[testDataRankSet]{Meta: &this.meta, Key: key, Data: data1}
	target.Initialize(context.Background(), nil, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaRankSet struct {
	table bool
}

func (this *testMetaRankSet) MajorKey(key any) string {
	return fmt.Sprintf("cmdrankset:%v", key)
}

func (this *testMetaRankSet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaRankSet) MinorTable() string {
	if this.table {
		return "cmdrankset"
	} // if

	return ""
}

type testDataRankSet struct {
	K     string `bson:"k"`
	Score int    `bson:"score"`
}
