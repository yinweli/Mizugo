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
	testdata.TestLeak
}

func (this *SuiteRand) SetupSuite() {
	this.Change("test-utils-rand")
}

func (this *SuiteRand) TearDownSuite() {
	this.Restore()
}

func (this *SuiteRand) TearDownTest() {
	this.GoLeak(this.T(), true)
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
