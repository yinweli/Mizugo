package nets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPSession(t *testing.T) {
	suite.Run(t, new(SuiteTCPSession))
}

type SuiteTCPSession struct {
	suite.Suite
	testdata.TestEnv
	ip      string
	port    int
	timeout time.Duration
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Change("test-tcpSession")
	this.ip = ""
	this.port = 3002
	this.timeout = time.Second
}

func (this *SuiteTCPSession) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPSession) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPSession) TestNewTCPSession() {
	assert.NotNil(this.T(), NewTCPSession(nil))
}

func (this *SuiteTCPSession) TestStartStop() {

}

func (this *SuiteTCPSession) TestSend() {

}

func (this *SuiteTCPSession) TestSessionID() {

}

func (this *SuiteTCPSession) TestRemoteAddr() {

}

func (this *SuiteTCPSession) TestLocalAddr() {

}
