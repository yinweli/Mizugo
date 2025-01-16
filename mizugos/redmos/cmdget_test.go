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

func TestCmdGet(t *testing.T) {
	suite.Run(t, new(SuiteCmdGet))
}

type SuiteCmdGet struct {
	suite.Suite
	trials.Catalog
	meta  metaGet
	major *Major
	minor *Minor
}

func (this *SuiteCmdGet) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdget"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdget")
}

func (this *SuiteCmdGet) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdGet) TestGet() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataGet{K: "redis+mongo", D: helps.RandStringDefault()}

	set := &Set[dataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataGet](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataGet](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target := &Get[dataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &Get[dataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: testdata.Unknown}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Data)

	target = &Get[dataGet]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Get[dataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &Get[dataGet]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaGet struct {
	table bool
}

func (this *metaGet) MajorKey(key any) string {
	return fmt.Sprintf("cmdget:%v", key)
}

func (this *metaGet) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaGet) MinorTable() string {
	if this.table {
		return "cmdget"
	} // if

	return ""
}

type dataGet struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
