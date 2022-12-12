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
	this.Change("test-nets-tcpListen")
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

func (this *SuiteTCPListen) TestListen() {
	testerl := newSessionTester()
	target := NewTCPListen(this.ip, this.port)
	go target.Listen(testerl.Complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Connect(testerc.Complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())
	assert.Nil(this.T(), target.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()

	verify := newSessionTester()
	target = NewTCPListen("!?", this.port)
	target.Listen(verify.Complete)
	assert.True(this.T(), verify.wait())
	assert.False(this.T(), verify.valid())

	verify = newSessionTester()
	target = NewTCPListen("192.168.0.1", this.port) // 故意要監聽錯誤位址才會引發錯誤
	target.Listen(verify.Complete)
	assert.True(this.T(), verify.wait())
	assert.False(this.T(), verify.valid())
}

func (this *SuiteTCPListen) TestTCPListen() {
	target := NewTCPListen(this.ip, this.port)
	assert.Equal(this.T(), net.JoinHostPort(this.ip, this.port), target.Address())
}
