package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestRedmomgr(t *testing.T) {
	suite.Run(t, new(SuiteRedmomgr))
}

type SuiteRedmomgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteRedmomgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-redmomgr"))
}

func (this *SuiteRedmomgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteRedmomgr) TestRedmomgr() {
	dbName := "dbName"
	majorName := "majorName"
	minorName := "minorName"
	mixedName := "mixedName"
	unknownName := "unknown"
	target := NewRedmomgr()
	assert.NotNil(this.T(), target)

	major, err := target.AddMajor(majorName, testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), major)
	_, err = target.AddMajor(majorName, testdata.RedisURI)
	assert.NotNil(this.T(), err)
	_, err = target.AddMajor(unknownName, testdata.RedisURIInvalid)
	assert.NotNil(this.T(), err)
	assert.NotNil(this.T(), target.GetMajor(majorName))
	assert.Nil(this.T(), target.GetMajor(unknownName))

	minor, err := target.AddMinor(minorName, testdata.MongoURI, dbName)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), minor)
	_, err = target.AddMinor(minorName, testdata.MongoURI, dbName)
	assert.NotNil(this.T(), err)
	_, err = target.AddMinor(minorName, testdata.MongoURI, "")
	assert.NotNil(this.T(), err)
	_, err = target.AddMinor(unknownName, testdata.MongoURIInvalid, dbName)
	assert.NotNil(this.T(), err)
	assert.NotNil(this.T(), target.GetMinor(minorName))
	assert.Nil(this.T(), target.GetMinor(unknownName))

	mixed, err := target.AddMixed(mixedName, majorName, minorName)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), mixed)
	_, err = target.AddMixed(mixedName, majorName, minorName)
	assert.NotNil(this.T(), err)
	_, err = target.AddMixed(unknownName, majorName, unknownName)
	assert.NotNil(this.T(), err)
	_, err = target.AddMixed(unknownName, unknownName, minorName)
	assert.NotNil(this.T(), err)
	assert.NotNil(this.T(), target.GetMixed(mixedName))
	assert.Nil(this.T(), target.GetMixed(unknownName))

	target.Finalize()
}
