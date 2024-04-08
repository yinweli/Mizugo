package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestCast(t *testing.T) {
	suite.Run(t, new(SuiteCast))
}

type SuiteCast struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteCast) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-cast")
}

func (this *SuiteCast) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteCast) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCast) TestCastPointer() {
	obj := &objTester1{}

	result, err := CastPointer[objTester1](obj)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), result)

	_, err = CastPointer[objTester2](obj)
	assert.NotNil(this.T(), err)

	_, err = CastPointer[objTester1](nil)
	assert.NotNil(this.T(), err)

	_, err = CastPointer[objTester2](nil)
	assert.NotNil(this.T(), err)
}

type objTester1 struct {
}

type objTester2 struct {
}
