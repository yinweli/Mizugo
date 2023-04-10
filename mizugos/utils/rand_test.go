package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestRand(t *testing.T) {
	suite.Run(t, new(SuiteRand))
}

type SuiteRand struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteRand) SetupSuite() {
	this.Env = testdata.EnvSetup("test-utils-rand")
}

func (this *SuiteRand) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteRand) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteRand) TestRandString() {
	letter := "0123456789"

	length := 32
	value := RandString(length, letter)
	assert.NotNil(this.T(), value)
	assert.Len(this.T(), value, length)
	fmt.Println(value)

	length = 64
	value = RandString(length, letter)
	assert.NotNil(this.T(), value)
	assert.Len(this.T(), value, length)
	fmt.Println(value)
}
