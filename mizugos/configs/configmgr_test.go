package configs

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestConfigmgr(t *testing.T) {
	suite.Run(t, new(SuiteConfigmgr))
}

type SuiteConfigmgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteConfigmgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-configs-configmgr"), testdata.PathEnv("configmgr"))
}

func (this *SuiteConfigmgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteConfigmgr) TestConfigmgr() {
	target := NewConfigmgr()
	this.NotNil(target)
	target.Reset()
}

func (this *SuiteConfigmgr) TestReadFile() {
	target := NewConfigmgr()
	target.AddPath(".")
	this.Nil(target.ReadFile("configmgr", "yaml"))
	this.Nil(target.ReadFile("configmgr", "yaml"))
	this.NotNil(target.ReadFile("configmgr", ""))
	this.NotNil(target.ReadFile(testdata.Unknown, "yaml"))
	this.Equal(1, target.Get("v1"))
}

func (this *SuiteConfigmgr) TestReadString() {
	target := NewConfigmgr()
	this.Nil(target.ReadString("v1: 1", "yaml"))
	this.Nil(target.ReadString("v2: 2", "yaml"))
	this.NotNil(target.ReadString("v2: 2", ""))
	this.NotNil(target.ReadString(testdata.Unknown, "yaml"))
	this.Equal(1, target.Get("v1"))
	this.Equal(2, target.Get("v2"))
}

func (this *SuiteConfigmgr) TestReadBuffer() {
	target := NewConfigmgr()
	this.Nil(target.ReadBuffer(bytes.NewBufferString("v1: 1"), "yaml"))
	this.Nil(target.ReadBuffer(bytes.NewBufferString("v2: 2"), "yaml"))
	this.NotNil(target.ReadBuffer(bytes.NewBufferString("v2: 2"), ""))
	this.NotNil(target.ReadBuffer(bytes.NewBufferString(testdata.Unknown), "yaml"))
	this.Equal(1, target.Get("v1"))
	this.Equal(2, target.Get("v2"))
}

func (this *SuiteConfigmgr) TestReadEnvironment() {
	_ = os.Setenv("ENV_V1", "1")
	target := NewConfigmgr()
	this.Nil(target.ReadEnvironment("env"))
	this.NotNil(target.ReadEnvironment(""))
	this.Equal("1", target.Get("v1"))
	_ = os.Unsetenv("ENV_V1")
}

func (this *SuiteConfigmgr) TestUnmarshal() {
	type data struct {
		V int
		S []string
	}

	actual := &data{}
	target := NewConfigmgr()
	target.AddPath(".")
	_ = target.ReadFile("configmgr", "yaml")
	this.Nil(target.Unmarshal("v2", &actual))
	this.NotNil(target.Unmarshal(testdata.Unknown, &actual))
	this.NotNil(target.Unmarshal("v2", nil))
	this.Equal(&data{
		V: 2,
		S: []string{"a", "b", "c"},
	}, actual)
}
