package redmos

import (
	"testing"

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
	target := NewRedmomgr()
	this.NotNil(target)
}

func (this *SuiteRedmomgr) TestMajor() {
	target := NewRedmomgr()
	major, err := target.AddMajor("major", testdata.RedisURI)
	this.Nil(err)
	this.NotNil(major)
	this.NotNil(target.GetMajor("major"))
	_, err = target.AddMajor("major", testdata.RedisURI)
	this.NotNil(err)
	_, err = target.AddMajor(testdata.Unknown, testdata.RedisURIInvalid)
	this.NotNil(err)
	this.Nil(target.GetMajor(testdata.Unknown))
	target.Finalize()
}

func (this *SuiteRedmomgr) TestMinor() {
	target := NewRedmomgr()
	minor, err := target.AddMinor("minor", testdata.MongoURI, "dbName")
	this.Nil(err)
	this.NotNil(minor)
	this.NotNil(target.GetMinor("minor"))
	_, err = target.AddMinor("minor", testdata.MongoURI, "dbName")
	this.NotNil(err)
	_, err = target.AddMinor(testdata.Unknown, testdata.MongoURI, "")
	this.NotNil(err)
	_, err = target.AddMinor(testdata.Unknown, testdata.MongoURIInvalid, "dbName")
	this.NotNil(err)
	this.Nil(target.GetMinor(testdata.Unknown))
	target.Finalize()
}

func (this *SuiteRedmomgr) TestMixed() {
	target := NewRedmomgr()
	_, _ = target.AddMajor("major", testdata.RedisURI)
	_, _ = target.AddMinor("minor", testdata.MongoURI, "dbName")
	mixed, err := target.AddMixed("mixed", "major", "minor")
	this.Nil(err)
	this.NotNil(mixed)
	this.NotNil(target.GetMixed("mixed"))
	_, err = target.AddMixed("mixed", testdata.Unknown, testdata.Unknown)
	this.NotNil(err)
	_, err = target.AddMixed(testdata.Unknown, "major", testdata.Unknown)
	this.NotNil(err)
	_, err = target.AddMixed(testdata.Unknown, testdata.Unknown, "minor")
	this.NotNil(err)
	this.Nil(target.GetMixed(testdata.Unknown))
	target.Finalize()
}
