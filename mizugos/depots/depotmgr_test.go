package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestDepotmgr(t *testing.T) {
	suite.Run(t, new(SuiteDepotmgr))
}

type SuiteDepotmgr struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteDepotmgr) SetupSuite() {
	this.Change("test-depots-depotmgr")
}

func (this *SuiteDepotmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteDepotmgr) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteDepotmgr) TestNewDepotmgr() {
	assert.NotNil(this.T(), NewDepotmgr())
}

func (this *SuiteDepotmgr) TestDepotmgr() {
	target := NewDepotmgr()
	majorName := "majorName"
	minorName := "minorName"
	mixedName := "mixedName"
	unknownName := "unknown"
	assert.Nil(this.T(), target.AddMajor(majorName, testdata.RedisURI))
	assert.NotNil(this.T(), target.AddMajor(majorName, ""))
	assert.NotNil(this.T(), target.AddMajor(unknownName, testdata.RedisURIInvalid))
	assert.NotNil(this.T(), target.GetMajor(majorName))
	assert.Nil(this.T(), target.GetMajor(unknownName))
	assert.Nil(this.T(), target.AddMinor(minorName, testdata.MongoURI))
	assert.NotNil(this.T(), target.AddMinor(minorName, ""))
	assert.NotNil(this.T(), target.AddMinor(unknownName, testdata.MongoURIInvalid))
	assert.NotNil(this.T(), target.GetMinor(minorName))
	assert.Nil(this.T(), target.GetMinor(unknownName))
	assert.Nil(this.T(), target.AddMixed(mixedName, majorName, minorName))
	assert.NotNil(this.T(), target.AddMixed(mixedName, majorName, minorName))
	assert.NotNil(this.T(), target.AddMixed(unknownName, majorName, unknownName))
	assert.NotNil(this.T(), target.AddMixed(unknownName, unknownName, minorName))
	assert.NotNil(this.T(), target.GetMixed(mixedName))
	assert.Nil(this.T(), target.GetMixed(unknownName))
	target.Stop()
}
