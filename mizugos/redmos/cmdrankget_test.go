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

func TestCmdRankGet(t *testing.T) {
	suite.Run(t, new(SuiteCmdRankGet))
}

type SuiteCmdRankGet struct {
	suite.Suite
	trials.Catalog
	meta  testMetaRankGet
	major *Major
	minor *Minor
}

func (this *SuiteCmdRankGet) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdrankget"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdrankget")
}

func (this *SuiteCmdRankGet) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdRankGet) TestRankGet() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataRankGet{K1: 100, K2: "1", D: "key1"}
	data2 := &testDataRankGet{K1: 200, K2: "2", D: "key2"}
	data3 := &testDataRankGet{K1: 150, K2: "3", D: "key3"}
	ahead := func(data *testDataRankGet) bson.M {
		return bson.M{"k1": bson.M{"$gt": data.K1}}
	}

	set := &Set[testDataRankGet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data1.D, Data: data1}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[testDataRankGet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data2.D, Data: data2}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[testDataRankGet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data3.D, Data: data3}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataRankGet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data1.D), data1))
	this.True(trials.MongoEqual[testDataRankGet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data2.D), data2))
	this.True(trials.MongoEqual[testDataRankGet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data3.D), data3))

	target := &RankGet[testDataRankGet]{Meta: &this.meta, Key: "key2", Ahead: ahead}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Equal(int64(1), target.Rank)

	target = &RankGet[testDataRankGet]{Meta: &this.meta, Key: "key1", Ahead: ahead}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Equal(int64(3), target.Rank)

	target = &RankGet[testDataRankGet]{Meta: &this.meta, Key: testdata.Unknown, Ahead: ahead}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.Equal(int64(0), target.Rank)

	target = &RankGet[testDataRankGet]{Meta: nil, Key: "k1", Ahead: ahead}
	this.NotNil(target.Prepare())

	target = &RankGet[testDataRankGet]{Meta: &this.meta, Key: "", Ahead: ahead}
	this.NotNil(target.Prepare())

	target = &RankGet[testDataRankGet]{Meta: &this.meta, Key: "k1", Ahead: nil}
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &RankGet[testDataRankGet]{Meta: &this.meta, Key: "k1", Ahead: ahead}
	this.NotNil(target.Prepare())
}

type testMetaRankGet struct {
	table bool
}

func (this *testMetaRankGet) MajorKey(key any) string {
	return fmt.Sprintf("cmdrankget:%v", key)
}

func (this *testMetaRankGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaRankGet) MinorTable() string {
	if this.table {
		return "cmdrankget"
	} // if

	return ""
}

type testDataRankGet struct {
	K1 int64  `bson:"k1"`
	K2 string `bson:"k2"`
	D  string `bson:"d"`
}
