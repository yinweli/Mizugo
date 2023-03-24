package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestRedmomgr(t *testing.T) {
	suite.Run(t, new(SuiteRedmomgr))
}

type SuiteRedmomgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteRedmomgr) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-redmomgr")
}

func (this *SuiteRedmomgr) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteRedmomgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteRedmomgr) TestNewRedmomgr() {
	assert.NotNil(this.T(), NewRedmomgr())
}

func (this *SuiteRedmomgr) TestRedmomgr() {
	target := NewRedmomgr()
	dbName := "dbName"
	majorName := "majorName"
	minorName := "minorName"
	mixedName := "mixedName"
	unknownName := "unknown"
	assert.Nil(this.T(), target.AddMajor(majorName, testdata.RedisURI, false))
	assert.NotNil(this.T(), target.AddMajor(majorName, "", false))
	assert.NotNil(this.T(), target.AddMajor(unknownName, testdata.RedisURIInvalid, false))
	assert.NotNil(this.T(), target.GetMajor(majorName))
	assert.Nil(this.T(), target.GetMajor(unknownName))
	assert.Nil(this.T(), target.AddMinor(minorName, testdata.MongoURI, dbName))
	assert.NotNil(this.T(), target.AddMinor(minorName, "", dbName))
	assert.NotNil(this.T(), target.AddMinor(minorName, testdata.MongoURI, ""))
	assert.NotNil(this.T(), target.AddMinor(unknownName, testdata.MongoURIInvalid, dbName))
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
