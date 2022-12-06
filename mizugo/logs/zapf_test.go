package logs

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestZapf(t *testing.T) {
	suite.Run(t, new(SuiteZapf))
}

type SuiteZapf struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteZapf) SetupSuite() {
	this.Change("test-logs-zapf")
}

func (this *SuiteZapf) TearDownSuite() {
	this.Restore()
}

func (this *SuiteZapf) TearDownTest() {
	goleak.VerifyNone(this.T())
}
