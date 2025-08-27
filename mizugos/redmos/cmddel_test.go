package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdDel(t *testing.T) {
	suite.Run(t, new(SuiteCmdDel))
}

type SuiteCmdDel struct {
	suite.Suite
	trials.Catalog
	meta  metaDel
	major *Major
	minor *Minor
}

func (this *SuiteCmdDel) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmddel"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmddel")
}

func (this *SuiteCmdDel) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdDel) TestDel() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataDel{K: "redis+mongo", D: helps.RandStringDefault()}

	set := &Set[dataDel]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataDel](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataDel](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target := &Del{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.False(this.T(), trials.RedisExist(this.major.Client(), this.meta.MajorKey(data.K)))
	assert.False(this.T(), trials.MongoExist(this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K)))

	target = &Del{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Del{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &Del{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaDel struct {
	table bool
}

func (this *metaDel) MajorKey(key any) string {
	return fmt.Sprintf("cmddel:%v", key)
}

func (this *metaDel) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaDel) MinorTable() string {
	if this.table {
		return "cmddel"
	} // if

	return ""
}

type dataDel struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
