package configs

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestConfigmgr(t *testing.T) {
	suite.Run(t, new(SuiteConfigmgr))
}

type SuiteConfigmgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteConfigmgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-configs-configmgr", "configmgr")
}

func (this *SuiteConfigmgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteConfigmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteConfigmgr) TestConfigmgr() {
	ext := "yaml"
	name1 := "configmgr"
	name2 := "!?"
	value1 := "valid: valid"
	value2 := "valid=valid"
	reader1 := bytes.NewBuffer([]byte(value1))
	reader2 := bytes.NewBuffer([]byte(value2))
	expected := configTester{
		Value1: 1,
		Value2: "2",
		Value3: []string{"a", "b", "c"},
	}
	valid := "valid"
	target := NewConfigmgr()
	assert.NotNil(this.T(), target)
	target.AddPath(".")
	assert.Nil(this.T(), target.ReadFile(name1, ext))
	assert.Nil(this.T(), target.ReadFile(name1, ext))
	assert.NotNil(this.T(), target.ReadFile(name2, ext))
	assert.Equal(this.T(), valid, target.Get(valid))
	target.Reset()

	target = NewConfigmgr()
	assert.Nil(this.T(), target.ReadString(value1, ext))
	assert.Nil(this.T(), target.ReadString(value1, ext))
	assert.NotNil(this.T(), target.ReadString(value2, ext))
	assert.Equal(this.T(), valid, target.Get(valid))
	target.Reset()

	target = NewConfigmgr()
	assert.Nil(this.T(), target.ReadBuffer(reader1, ext))
	assert.Nil(this.T(), target.ReadBuffer(reader1, ext))
	assert.NotNil(this.T(), target.ReadBuffer(reader2, ext))
	assert.Equal(this.T(), valid, target.Get(valid))
	target.Reset()

	target = NewConfigmgr()
	target.AddPath(".")
	assert.Nil(this.T(), target.ReadFile(name1, ext))
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

	target = NewConfigmgr()
	actual := configTester{}
	target.AddPath(".")
	assert.Nil(this.T(), target.ReadFile(name1, ext))
	assert.Nil(this.T(), target.Unmarshal("object", &actual))
	assert.NotNil(this.T(), target.Unmarshal("!?", &actual))
	assert.NotNil(this.T(), target.Unmarshal("object", nil))
	assert.Equal(this.T(), expected, actual)
	target.Reset()
}

type configTester struct {
	Value1 int
	Value2 string
	Value3 []string
}
