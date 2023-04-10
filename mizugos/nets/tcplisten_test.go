package nets

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPListen(t *testing.T) {
	suite.Run(t, new(SuiteTCPListen))
}

type SuiteTCPListen struct {
	suite.Suite
	testdata.Env
	host host
}

func (this *SuiteTCPListen) SetupSuite() {
	this.Env = testdata.EnvSetup("test-nets-tcpListen")
	this.host = host{port: "9001"}
}

func (this *SuiteTCPListen) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteTCPListen) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteTCPListen) TestNewTCPListen() {
	assert.NotNil(this.T(), NewTCPListen(this.host.ip, this.host.port))
}

func (this *SuiteTCPListen) TestListen() {
	testl := newTester(true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testc.get().Stop()
	assert.Nil(this.T(), listen.Stop())

	testl = newTester(true, true, true)
	listen = NewTCPListen("!?", this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testl.valid())

	testl = newTester(true, true, true)
	listen = NewTCPListen("192.168.0.1", this.host.port) // 故意要接聽錯誤位址才會引發錯誤
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testl.valid())
}

func (this *SuiteTCPListen) TestAddress() {
	target := NewTCPListen(this.host.ip, this.host.port)
	assert.Equal(this.T(), net.JoinHostPort(this.host.ip, this.host.port), target.Address())
}
