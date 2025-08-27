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

func TestCmdGetEx(t *testing.T) {
	suite.Run(t, new(SuiteCmdGetEx))
}

type SuiteCmdGetEx struct {
	suite.Suite
	trials.Catalog
	meta  metaGetEx
	major *Major
	minor *Minor
}

func (this *SuiteCmdGetEx) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdgetex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdgetex")
}

func (this *SuiteCmdGetEx) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdGetEx) TestGetEx() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataGetEx{K: "redis+mongo", D: helps.RandStringDefault()}

	set := &Set[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K, Data: data}
	set.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), set.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), set.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), trials.RedisCompare[dataGetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataGetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))

	target := &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: false, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &GetEx[dataGetEx]{MajorEnable: false, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.True(this.T(), cmp.Equal(data, target.Data))

	target = &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: testdata.Unknown}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Data)

	target = &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: nil, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	target = &GetEx[dataGetEx]{MajorEnable: true, MinorEnable: true, Meta: &this.meta, Expire: trials.Timeout * 2, Key: data.K}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	trials.WaitTimeout()
	assert.True(this.T(), trials.RedisCompare[dataGetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataGetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
	trials.WaitTimeout()
	assert.False(this.T(), trials.RedisCompare[dataGetEx](this.major.Client(), this.meta.MajorKey(data.K), data))
	assert.True(this.T(), trials.MongoCompare[dataGetEx](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(data.K), data))
}

type metaGetEx struct {
	table bool
}

func (this *metaGetEx) MajorKey(key any) string {
	return fmt.Sprintf("cmdgetex:%v", key)
}

func (this *metaGetEx) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaGetEx) MinorTable() string {
	if this.table {
		return "cmdgetex"
	} // if

	return ""
}

type dataGetEx struct {
	K string `bson:"k"`
	D string `bson:"d"`
}
