package redmos

import (
	"context"
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
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()

	target := &Index{Name: "index", Table: "cmdindex", Field: "field", Order: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())

	target = &Index{Name: "index", Table: "cmdindex", Field: "field", Order: 1}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())

	target = &Index{Name: "", Table: "cmdindex", Field: "field", Order: 1}
	assert.NotNil(this.T(), target.Prepare())

	target = &Index{Name: "index", Table: "", Field: "field", Order: 1}
	assert.NotNil(this.T(), target.Prepare())

	target = &Index{Name: "index", Table: "cmdindex", Field: "", Order: 1}
	assert.NotNil(this.T(), target.Prepare())

	target = &Index{Name: "index", Table: "cmdindex", Field: "field", Order: 0}
	assert.NotNil(this.T(), target.Prepare())
}
