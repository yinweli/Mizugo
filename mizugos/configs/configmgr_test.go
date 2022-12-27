package configs

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestConfigmgr(t *testing.T) {
	suite.Run(t, new(SuiteConfigmgr))
}

type SuiteConfigmgr struct {
	suite.Suite
	testdata.TestEnv
	name1   string
	name2   string
	value1  string
	value2  string
	reader1 io.Reader
	reader2 io.Reader
	object  configTester
	ext     string
	valid   string
}

func (this *SuiteConfigmgr) SetupSuite() {
	this.Change("test-configs-configmgr")
	this.name1 = "configmgr"
	this.name2 = "!?"
	this.value1 = "valid: valid"
	this.value2 = "valid=valid"
	this.reader1 = bytes.NewBuffer([]byte(this.value1))
	this.reader2 = bytes.NewBuffer([]byte(this.value2))
	this.object = configTester{
		Value1: 1,
		Value2: "2",
		Value3: []string{"a", "b", "c"},
	}
	this.ext = "yaml"
	this.valid = "valid"
}

func (this *SuiteConfigmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteConfigmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteConfigmgr) TestNewConfigmgr() {
	assert.NotNil(this.T(), NewConfigmgr())
}

func (this *SuiteConfigmgr) TestReadFile() {
	target := NewConfigmgr()
	target.AddPath(".")
	assert.Nil(this.T(), target.ReadFile(this.name1, this.ext))
	assert.Nil(this.T(), target.ReadFile(this.name1, this.ext))
	assert.NotNil(this.T(), target.ReadFile(this.name2, this.ext))
	assert.Equal(this.T(), this.valid, target.Get(this.valid))
	target.Reset()
}

func (this *SuiteConfigmgr) TestReadString() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadString(this.value1, this.ext))
	assert.Nil(this.T(), target.ReadString(this.value1, this.ext))
	assert.NotNil(this.T(), target.ReadString(this.value2, this.ext))
	assert.Equal(this.T(), this.valid, target.Get(this.valid))
	target.Reset()
}

func (this *SuiteConfigmgr) TestReadBuffer() {
	target := NewConfigmgr()
	assert.Nil(this.T(), target.ReadBuffer(this.reader1, this.ext))
	assert.Nil(this.T(), target.ReadBuffer(this.reader1, this.ext))
	assert.NotNil(this.T(), target.ReadBuffer(this.reader2, this.ext))
	assert.Equal(this.T(), this.valid, target.Get(this.valid))
	target.Reset()
}

func (this *SuiteConfigmgr) TestGet() {
	target := NewConfigmgr()
	target.AddPath(".")
	assert.Nil(this.T(), target.ReadFile(this.name1, this.ext))
	assert.Equal(this.T(), true, target.GetBool("valueb"))
	assert.Equal(this.T(), 1, target.GetInt("valuei"))
	assert.Equal(this.T(), int32(100000000), target.GetInt32("valuei32"))
	assert.Equal(this.T(), int64(100000000000), target.GetInt64("valuei64"))
	assert.Equal(this.T(), uint(2), target.GetUInt("valueu"))
	assert.Equal(this.T(), uint32(200000000), target.GetUInt32("valueu32"))
	assert.Equal(this.T(), uint64(200000000000), target.GetUInt64("valueu64"))
	assert.Equal(this.T(), 3.33, target.GetFloat("valuef"))
	assert.Equal(this.T(), "string", target.GetString("values"))
	assert.Equal(this.T(), []int{1, 2, 3}, target.GetIntSlice("valuelisti"))
	assert.Equal(this.T(), []string{"a", "b", "c"}, target.GetStringSlice("valuelists"))
	assert.Equal(this.T(), "2020-12-31 09:30:30", target.GetTime("valuet").Format("2006-01-02 15:04:05"))
	assert.Equal(this.T(), time.Second*360, target.GetDuration("valued"))
	assert.Equal(this.T(), uint(1024), target.GetSizeInBytes("valuec"))
	target.Reset()
}

func (this *SuiteConfigmgr) TestUnmarshal() {
	target := NewConfigmgr()
	object := configTester{}
	target.AddPath(".")
	assert.Nil(this.T(), target.ReadFile(this.name1, this.ext))
	assert.Nil(this.T(), target.Unmarshal("object", &object))
	assert.NotNil(this.T(), target.Unmarshal("!?", &object))
	assert.NotNil(this.T(), target.Unmarshal("object", nil))
	assert.Equal(this.T(), this.object, object)
	target.Reset()
}

type configTester struct {
	Value1 int
	Value2 string
	Value3 []string
}
