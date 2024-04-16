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
	test := newTester(true, true, true)
	target := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	assert.NotNil(this.T(), target)
	target.Connect(test.bind, test.unbind, test.wrong)

	time.Sleep(trials.Timeout)
	assert.True(this.T(), test.valid())
	test.get().StopWait()

	test = newTester(true, true, true)
	target = NewTCPConnect("!?", addr.port, trials.Timeout)
	target.Connect(test.bind, test.unbind, test.wrong)

	time.Sleep(trials.Timeout)
	assert.False(this.T(), test.valid())

	test = newTester(true, true, true)
	target = NewTCPConnect(addr.ip, "9999", trials.Timeout) // 故意連線到不開放的埠號才會引發錯誤
	target.Connect(test.bind, test.unbind, test.wrong)

	time.Sleep(trials.Timeout * 2) // 因為錯誤會是timeout, 所以要等待長一點
	assert.False(this.T(), test.valid())

	target = NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	assert.Equal(this.T(), net.JoinHostPort(addr.ip, addr.port), target.Address())
}
