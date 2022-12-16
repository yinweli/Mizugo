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

func TestTCPListen(t *testing.T) {
	suite.Run(t, new(SuiteTCPListen))
}

type SuiteTCPListen struct {
	suite.Suite
	testdata.TestEnv
	host host
}

func (this *SuiteTCPListen) SetupSuite() {
	this.Change("test-nets-tcpListen")
	this.host = host{ip: "", port: "3001"}
}

func (this *SuiteTCPListen) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPListen) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPListen) TestNewTCPListen() {
	assert.NotNil(this.T(), NewTCPListen(this.host.ip, this.host.port))
}

func (this *SuiteTCPListen) TestListen() {
	testerl := newCompleteTester()
	target := NewTCPListen(this.host.ip, this.host.port)
	target.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testerc)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())
	assert.Nil(this.T(), target.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()

	tester := newCompleteTester()
	target = NewTCPListen("!?", this.host.port)
	target.Listen(tester)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), tester.valid())

	tester = newCompleteTester()
	target = NewTCPListen("192.168.0.1", this.host.port) // 故意要接聽錯誤位址才會引發錯誤
	target.Listen(tester)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), tester.valid())
}

func (this *SuiteTCPListen) TestAddress() {
	target := NewTCPListen(this.host.ip, this.host.port)
	assert.Equal(this.T(), net.JoinHostPort(this.host.ip, this.host.port), target.Address())
}
