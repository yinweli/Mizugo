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
	testdata.TestEnv
}

func (this *SuiteRand) SetupSuite() {
	this.TBegin("test-utils-rand", "")
}

func (this *SuiteRand) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteRand) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteRand) TestRandString() {
	length := 32
	value := RandString(length)
	assert.NotNil(this.T(), value)
	assert.Len(this.T(), value, length)
	fmt.Println(value)

	length = 64
	value = RandString(length)
	assert.NotNil(this.T(), value)
	assert.Len(this.T(), value, length)
	fmt.Println(value)
}
