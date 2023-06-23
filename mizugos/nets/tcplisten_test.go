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
}

func (this *SuiteTCPListen) SetupSuite() {
	this.Env = testdata.EnvSetup("test-nets-tcpListen")
}

func (this *SuiteTCPListen) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteTCPListen) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteTCPListen) TestTCPListen() {
	addr := host{port: "9001"}
	testl := newTester(true, true, true)
	target := NewTCPListen(addr.ip, addr.port)
	assert.NotNil(this.T(), target)
	target.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testc.get().Stop()
	assert.Nil(this.T(), target.Stop())

	testl = newTester(true, true, true)
	target = NewTCPListen("!?", addr.port)
	target.Listen(testl.bind, testl.unbind, testl.wrong)

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testl.valid())

	testl = newTester(true, true, true)
	target = NewTCPListen("192.168.0.1", addr.port) // 故意要接聽錯誤位址才會引發錯誤
	target.Listen(testl.bind, testl.unbind, testl.wrong)

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testl.valid())

	target = NewTCPListen(addr.ip, addr.port)
	assert.Equal(this.T(), net.JoinHostPort(addr.ip, addr.port), target.Address())
}
