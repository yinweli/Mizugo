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
	ip      string
	port    string
	timeout time.Duration
}

func (this *SuiteTCPListen) SetupSuite() {
	this.Change("test-tcpListen")
	this.ip = ""
	this.port = "3001"
	this.timeout = time.Second
}

func (this *SuiteTCPListen) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPListen) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPListen) TestNewTCPListen() {
	assert.NotNil(this.T(), NewTCPListen(this.ip, this.port))
}

func (this *SuiteTCPListen) TestStartStop() {
	sessionl := newTestSession("listen server", this.timeout)
	target := NewTCPListen(this.ip, this.port)
	go target.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("listen client", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.True(this.T(), sessionc.Wait())
	assert.Nil(this.T(), target.Stop())

	session := newTestSession("listen error ip1", this.timeout)
	target = NewTCPListen("!?", this.port)
	target.Start(session.Complete)
	assert.True(this.T(), session.Wait())
	assert.NotNil(this.T(), session.Error())

	session = newTestSession("listen error ip2", this.timeout)
	target = NewTCPListen("192.168.0.1", this.port) // 故意要監聽錯誤位址才會引發錯誤
	target.Start(session.Complete)
	assert.True(this.T(), session.Wait())
	assert.NotNil(this.T(), session.Error())
}

func (this *SuiteTCPListen) TestAddress() {
	target := NewTCPListen(this.ip, this.port)
	assert.Equal(this.T(), net.JoinHostPort(this.ip, this.port), target.Address())
}
