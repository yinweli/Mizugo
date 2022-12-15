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
	host    host
	timeout time.Duration
}

func (this *SuiteTCPConnect) SetupSuite() {
	this.Change("test-nets-tcpConnect")
	this.host = host{ip: "google.com", port: "80"}
	this.timeout = time.Second
}

func (this *SuiteTCPConnect) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPConnect) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPConnect) TestNewTCPConnect() {
	assert.NotNil(this.T(), NewTCPConnect(this.host.ip, this.host.port, this.timeout))
}

func (this *SuiteTCPConnect) TestConnect() {
	tester := newCompleteTester()
	target := NewTCPConnect(this.host.ip, this.host.port, this.timeout)
	target.Connect(tester)
	time.Sleep(this.timeout)
	assert.True(this.T(), tester.valid())
	tester.get().StopWait()

	tester = newCompleteTester()
	target = NewTCPConnect("!?", this.host.port, this.timeout)
	target.Connect(tester)
	time.Sleep(this.timeout)
	assert.False(this.T(), tester.valid())

	tester = newCompleteTester()
	target = NewTCPConnect(this.host.ip, "9999", this.timeout) // 故意連線到不開放的埠號才會引發錯誤
	target.Connect(tester)
	time.Sleep(this.timeout)
	assert.False(this.T(), tester.valid())
}

func (this *SuiteTCPConnect) TestAddress() {
	target := NewTCPConnect(this.host.ip, this.host.port, this.timeout)
	assert.Equal(this.T(), net.JoinHostPort(this.host.ip, this.host.port), target.Address())
}
