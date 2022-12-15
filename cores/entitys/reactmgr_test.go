package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestReactmgr(t *testing.T) {
	suite.Run(t, new(SuiteReactmgr))
}

type SuiteReactmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteReactmgr) SetupSuite() {
	this.Change("test-entitys-reactmgr")
}

func (this *SuiteReactmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteReactmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteReactmgr) TestNewReactmgr() {
	assert.NotNil(this.T(), NewReactmgr())
}

func (this *SuiteReactmgr) TestSetGet() {
	react := newReactTester()
	target := NewReactmgr()
	target.Set(react)
	assert.Equal(this.T(), react, target.Get())
}
