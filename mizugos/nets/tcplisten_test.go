package nets

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPListen(t *testing.T) {
	suite.Run(t, new(SuiteTCPListen))
}

type SuiteTCPListen struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteTCPListen) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-nets-tcpListen"))
}

func (this *SuiteTCPListen) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteTCPListen) TestTCPListen() {
	addr := host{port: "9001"}
	testl := newTester(true, true, true)
	target := NewTCPListen(addr.ip, addr.port)
	assert.NotNil(this.T(), target)
	target.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(trials.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(trials.Timeout)
	testc.get().Stop()
	assert.Nil(this.T(), target.Stop())

	testl = newTester(true, true, true)
	target = NewTCPListen("!?", addr.port)
	target.Listen(testl.bind, testl.unbind, testl.wrong)

	time.Sleep(trials.Timeout)
	assert.False(this.T(), testl.valid())

	testl = newTester(true, true, true)
	target = NewTCPListen("192.168.0.1", addr.port) // 故意要接聽錯誤位址才會引發錯誤
	target.Listen(testl.bind, testl.unbind, testl.wrong)

	time.Sleep(trials.Timeout)
	assert.False(this.T(), testl.valid())

	target = NewTCPListen(addr.ip, addr.port)
	assert.Equal(this.T(), net.JoinHostPort(addr.ip, addr.port), target.Address())
}
