package redmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdIndex(t *testing.T) {
	suite.Run(t, new(SuiteCmdIndex))
}

type SuiteCmdIndex struct {
	suite.Suite
	trials.Catalog
	meta  metaIndex
	major *Major
	minor *Minor
}

func (this *SuiteCmdIndex) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdindex"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdindex")
}

func (this *SuiteCmdIndex) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdIndex) TestIndex() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()

	target := &Index{Meta: &this.meta, Name: "field", Order: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())

	target = &Index{Meta: &this.meta, Name: "field", Order: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())

	target = &Index{Meta: nil, Name: "field", Order: 1}
	assert.NotNil(this.T(), target.Prepare())

	target = &Index{Meta: &this.meta, Name: "", Order: 1}
	assert.NotNil(this.T(), target.Prepare())

	target = &Index{Meta: &this.meta, Name: "field", Order: 0}
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &Index{Meta: &this.meta, Name: "field", Order: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &Index{Meta: &this.meta, Name: "field", Order: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

type metaIndex struct {
	table bool
	field bool
}

func (this *metaIndex) MajorKey(key any) string {
	return fmt.Sprintf("cmdindex:%v", key)
}

func (this *metaIndex) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaIndex) MinorTable() string {
	if this.table {
		return "cmdindex"
	} // if

	return ""
}

func (this *metaIndex) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}
