package nets

import (
	"net"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestTCPConnect(t *testing.T) {
	suite.Run(t, new(SuiteTCPConnect))
}

type SuiteTCPConnect struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteTCPConnect) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-nets-tcpConnect"))
}

func (this *SuiteTCPConnect) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteTCPConnect) TestTCPConnect() {
	addr := host{ip: "google.com", port: "80"}
	test := newTestNet(true, true, true)
	target := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	this.NotNil(target)
	target.Connect(test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	this.True(test.Valid())
	test.Get().StopWait()

	test = newTestNet(true, true, true)
	target = NewTCPConnect("!?", addr.port, trials.Timeout)
	target.Connect(test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	this.False(test.Valid())

	test = newTestNet(true, true, true)
	target = NewTCPConnect(addr.ip, "9999", trials.Timeout) // 故意連線到不開放的埠號才會引發錯誤
	target.Connect(test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout(trials.Timeout * 2) // 因為錯誤會是timeout, 所以要等待長一點
	this.False(test.Valid())

	target = NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	this.Equal(net.JoinHostPort(addr.ip, addr.port), target.Address())
}
