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
	donel := newDoneTester()
	target := NewTCPListen(this.host.ip, this.host.port)
	target.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())
	assert.Nil(this.T(), target.Stop())
	donel.get().StopWait()
	donec.get().StopWait()

	done := newDoneTester()
	target = NewTCPListen("!?", this.host.port)
	target.Listen(done.done)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), done.valid())

	done = newDoneTester()
	target = NewTCPListen("192.168.0.1", this.host.port) // 故意要接聽錯誤位址才會引發錯誤
	target.Listen(done.done)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), done.valid())
}

func (this *SuiteTCPListen) TestAddress() {
	target := NewTCPListen(this.host.ip, this.host.port)
	assert.Equal(this.T(), net.JoinHostPort(this.host.ip, this.host.port), target.Address())
}
