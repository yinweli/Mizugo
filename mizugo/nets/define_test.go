package nets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestDefine(t *testing.T) {
	suite.Run(t, new(SuiteDefine))
}

type SuiteDefine struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteDefine) SetupSuite() {
	this.Change("test-nets-tcpSession")
}

func (this *SuiteDefine) TearDownSuite() {
	this.Restore()
}

func (this *SuiteDefine) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteDefine) TestComplete() {
	valid := false
	complete := Complete(func(_ Sessioner, _ error) {
		valid = true
	})
	complete.Complete(nil, nil)
	assert.True(this.T(), valid)
}

func (this *SuiteDefine) TestRelease() {
	valid := false
	release := Release(func() {
		valid = true
	})
	release.Release()
	assert.True(this.T(), valid)
}
