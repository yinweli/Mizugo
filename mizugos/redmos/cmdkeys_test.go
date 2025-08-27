package redmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestCmdKeys(t *testing.T) {
	suite.Run(t, new(SuiteCmdKeys))
}

type SuiteCmdKeys struct {
	suite.Suite
	trials.Catalog
	major *Major
	minor *Minor
}

func (this *SuiteCmdKeys) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdkeys"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdkeys")
}

func (this *SuiteCmdKeys) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdKeys) TestKeys() {
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	this.major.Client().Set(context.Background(), "key1", "value0", 0)
	this.major.Client().Set(context.Background(), "key2", "value1", 0)

	target := &Keys{Pattern: "key*"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	assert.ElementsMatch(this.T(), target.Data, []string{"key1", "key2"})

	target = &Keys{Pattern: "a*"}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	assert.Empty(this.T(), target.Data)

	target = &Keys{Pattern: ""}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}
