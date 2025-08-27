package helps

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestReflect(t *testing.T) {
	suite.Run(t, new(SuiteReflect))
}

type SuiteReflect struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteReflect) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-reflect"))
}

func (this *SuiteReflect) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteReflect) TestReflect() {
	target := reflect.ValueOf(&testReflect{Value: testdata.Unknown})
	result, ok := ReflectFieldValue[string](target, "Value")
	assert.True(this.T(), ok)
	assert.Equal(this.T(), testdata.Unknown, result)
	_, ok = ReflectFieldValue[string](target, testdata.Unknown)
	assert.False(this.T(), ok)
}

type testReflect struct {
	Value string
}
