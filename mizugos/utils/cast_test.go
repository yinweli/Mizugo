package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestCast(t *testing.T) {
	suite.Run(t, new(SuiteCast))
}

type SuiteCast struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteCast) SetupSuite() {
	this.Change("test-utils-cast")
}

func (this *SuiteCast) TearDownSuite() {
	this.Restore()
}

func (this *SuiteCast) TearDownTest() {
	goleak.VerifyNone(this.T())
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
