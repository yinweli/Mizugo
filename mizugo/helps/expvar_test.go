package helps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestExpvar(t *testing.T) {
	suite.Run(t, new(SuiteExpvar))
}

type SuiteExpvar struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteExpvar) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-expvar")
}

func (this *SuiteExpvar) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteExpvar) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteExpvar) TestExpvarStr() {
	expvarStat := []ExpvarStat{
		{"value1", nil},
		{"value2", 1001},
		{"value3", "data"},
		{"value4", time.Second},
		{"value5", struct {
			Data1 int
			Data2 string
		}{1, "a"}},
	}
	expected := "{\"value1\": \"<nil>\", \"value2\": 1001, \"value3\": \"data\", \"value4\": \"1s\", \"value5\": \"{1 a}\"}"
	assert.Equal(this.T(), expected, ExpvarStr(expvarStat))
}

func (this *SuiteExpvar) TestExpvarStat() {
	assert.True(this.T(), ExpvarStat{Name: "", Data: nil}.stringType())
	assert.False(this.T(), ExpvarStat{Name: "", Data: 1001}.stringType())
	assert.True(this.T(), ExpvarStat{Name: "", Data: "data"}.stringType())
	assert.True(this.T(), ExpvarStat{Name: "", Data: time.Second}.stringType())
	assert.True(this.T(), ExpvarStat{
		Name: "",
		Data: struct {
			Data1 int
			Data2 string
		}{1, "a"},
	}.stringType())
}
