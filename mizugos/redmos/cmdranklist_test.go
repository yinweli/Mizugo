package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdRankList(t *testing.T) {
	suite.Run(t, new(SuiteCmdRankList))
}

type SuiteCmdRankList struct {
	suite.Suite
	trials.Catalog
	meta  testMetaRankList
	major *Major
	minor *Minor
}

func (this *SuiteCmdRankList) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdranklist"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdranklist")
}

func (this *SuiteCmdRankList) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdRankList) TestRankList() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataRankList{K1: 100, K2: "1", D: "key1"}
	data2 := &testDataRankList{K1: 200, K2: "2", D: "key2"}
	data3 := &testDataRankList{K1: 150, K2: "3", D: "key3"}

	set := &Set[testDataRankList]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data1.D, Data: data1}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[testDataRankList]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data2.D, Data: data2}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[testDataRankList]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data3.D, Data: data3}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataRankList](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.D), data1))
	this.True(trials.MongoEqual[testDataRankList](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data2.D), data2))
	this.True(trials.MongoEqual[testDataRankList](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data3.D), data3))

	target := &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Sort: []SortField{{Field: "k1", Order: -1}}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Len(target.Data, 3)
	this.Equal(int64(200), target.Data[0].K1)

	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 2, Sort: []SortField{{Field: "k1", Order: -1}}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Len(target.Data, 2)

	target = &RankList[testDataRankList]{Meta: nil, Limit: 3, Sort: []SortField{{Field: "k1", Order: -1}}}
	this.NotNil(target.Prepare())

	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 0, Sort: []SortField{{Field: "k1", Order: -1}}}
	this.NotNil(target.Prepare())

	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Sort: []SortField{}}
	this.NotNil(target.Prepare())

	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Sort: []SortField{{Field: "", Order: -1}}}
	this.NotNil(target.Prepare())

	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Sort: []SortField{{Field: "k1", Order: 0}}}
	this.NotNil(target.Prepare())

	this.meta.table = true
	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Filter: bson.D{{Key: "k2", Value: "2"}}, Sort: []SortField{{Field: "k1", Order: -1}}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Len(target.Data, 1)
	this.Equal(int64(200), target.Data[0].K1)

	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Filter: bson.D{}, Sort: []SortField{{Field: "k1", Order: -1}}}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Len(target.Data, 3)

	this.meta.table = false
	target = &RankList[testDataRankList]{Meta: &this.meta, Limit: 3, Sort: []SortField{{Field: "k1", Order: -1}}}
	this.NotNil(target.Prepare())
}

type testMetaRankList struct {
	table bool
}

func (this *testMetaRankList) MajorKey(key any) string {
	return fmt.Sprintf("cmdranklist:%v", key)
}

func (this *testMetaRankList) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaRankList) MinorTable() string {
	if this.table {
		return "cmdranklist"
	} // if

	return ""
}

type testDataRankList struct {
	K1 int64  `bson:"k1"`
	K2 string `bson:"k2"`
	D  string `bson:"d"`
}
