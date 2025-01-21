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

func TestCmdFind(t *testing.T) {
	suite.Run(t, new(SuiteCmdFind))
}

type SuiteCmdFind struct {
	suite.Suite
	trials.Catalog
	meta  metaFind
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
	data1 := &dataFind{K: "key1", D: helps.RandStringDefault()}
	data2 := &dataFind{K: "key2", D: helps.RandStringDefault()}

	set := &Set[dataFind]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data1.K, Data: data1}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	set = &Set[dataFind]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data2.K, Data: data2}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.MongoCompare[dataFind](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey("key1"), data1))
	assert.True(this.T(), trials.MongoCompare[dataFind](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey("key2"), data2))

	target := &Find{Meta: &this.meta, Pattern: "^key"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	assert.ElementsMatch(this.T(), target.Data, []string{"key1", "key2"})

	target = &Find{Meta: &this.meta, Pattern: "[abc"} // 故意寫錯誤的正則表達式
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.NotNil(this.T(), target.Complete())

	target = &Find{Meta: nil, Pattern: "^key"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Find{Meta: &this.meta, Pattern: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &Find{Meta: &this.meta, Pattern: "^key"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaFind struct {
	table bool
}

func (this *metaFind) MajorKey(key any) string {
	return fmt.Sprintf("cmdfind:%v", key)
}

func (this *metaFind) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaFind) MinorTable() string {
	if this.table {
		return "cmdfind"
	} // if

	return ""
}

type dataFind struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
