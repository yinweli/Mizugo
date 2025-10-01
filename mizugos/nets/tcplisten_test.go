package nets

import (
	"net"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	testl := newTestNet(true, true, true)
	target := NewTCPListen(addr.ip, addr.port)
	this.NotNil(target)
	target.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTestNet(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testc.Valid())

	trials.WaitTimeout()
	testc.Get().Stop()
	this.Nil(target.Stop())

	testl = newTestNet(true, true, true)
	target = NewTCPListen("!?", addr.port)
	target.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	trials.WaitTimeout()
	this.False(testl.Valid())

	testl = newTestNet(true, true, true)
	target = NewTCPListen("192.168.0.1", addr.port) // 故意要接聽錯誤位址才會引發錯誤
	target.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	trials.WaitTimeout()
	this.False(testl.Valid())

	target = NewTCPListen(addr.ip, addr.port)
	this.Equal(net.JoinHostPort(addr.ip, addr.port), target.Address())
}
