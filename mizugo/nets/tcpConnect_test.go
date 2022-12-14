package nets

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPConnect(t *testing.T) {
	suite.Run(t, new(SuiteTCPConnect))
}

type SuiteTCPConnect struct {
	suite.Suite
	testdata.TestEnv
	ip      string
	port    string
	timeout time.Duration
}

func (this *SuiteTCPConnect) SetupSuite() {
	this.Change("test-nets-tcpConnect")
	this.ip = "google.com"
	this.port = "80"
	this.timeout = time.Second
}

func (this *SuiteTCPConnect) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPConnect) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPConnect) TestNewTCPConnect() {
	assert.NotNil(this.T(), NewTCPConnect(this.ip, this.port, this.timeout))
}

func (this *SuiteTCPConnect) TestConnect() {
	tester := newCompleteTester()
	target := NewTCPConnect(this.ip, this.port, this.timeout)
	target.Connect(tester)
	time.Sleep(this.timeout)
	assert.True(this.T(), tester.valid())
	tester.get().StopWait()

	tester = newCompleteTester()
	target = NewTCPConnect("!?", this.port, this.timeout)
	target.Connect(tester)
	time.Sleep(this.timeout)
	assert.False(this.T(), tester.valid())

	tester = newCompleteTester()
	target = NewTCPConnect(this.ip, "3000", this.timeout) // 故意連線到不開放的埠號才會引發錯誤
	target.Connect(tester)
	time.Sleep(this.timeout)
	assert.False(this.T(), tester.valid())
}

func (this *SuiteTCPConnect) TestTCPConnect() {
	target := NewTCPConnect(this.ip, this.port, this.timeout)
	assert.Equal(this.T(), net.JoinHostPort(this.ip, this.port), target.Address())
}
