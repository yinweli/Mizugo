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

func TestCmdFind(t *testing.T) {
	suite.Run(t, new(SuiteCmdFind))
}

type SuiteCmdFind struct {
	suite.Suite
	trials.Catalog
	meta  testMetaFind
	major *Major
	minor *Minor
}

func (this *SuiteCmdFind) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdfind"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdfind")
}

func (this *SuiteCmdFind) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdFind) TestFind() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data1 := &testDataFind{K: "key1", D: helps.RandStringDefault()}
	data2 := &testDataFind{K: "key2", D: helps.RandStringDefault()}

	set := &Set[testDataFind]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[testDataFind]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(set.Complete())
	_ = minorSubmit.Exec(context.Background())
	this.True(trials.MongoEqual[testDataFind](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey("key1"), data1))
	this.True(trials.MongoEqual[testDataFind](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey("key2"), data2))

	target := &Find{Meta: &this.meta, Pattern: "^key"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(target.Complete())
	this.ElementsMatch(target.Data, []string{"key1", "key2"})

	target = &Find{Meta: &this.meta, Pattern: "[abc"} // 故意寫錯誤的正則表達式
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.Nil(target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.NotNil(target.Complete())

	target = &Find{Meta: nil, Pattern: "^key"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	target = &Find{Meta: &this.meta, Pattern: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())

	this.meta.table = false
	target = &Find{Meta: &this.meta, Pattern: "^key"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	this.NotNil(target.Prepare())
}

type testMetaFind struct {
	table bool
}

func (this *testMetaFind) MajorKey(key any) string {
	return fmt.Sprintf("cmdfind:%v", key)
}

func (this *testMetaFind) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *testMetaFind) MinorTable() string {
	if this.table {
		return "cmdfind"
	} // if

	return ""
}

type testDataFind struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
