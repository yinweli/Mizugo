package nets

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPConnect(t *testing.T) {
	suite.Run(t, new(SuiteTCPConnect))
}

type SuiteTCPConnect struct {
	suite.Suite
	testdata.TestEnv
	ip      string
	port    string
	timeout time.Duration
}

func (this *SuiteTCPConnect) SetupSuite() {
	this.Change("test-tcpConnect")
	this.ip = "google.com"
	this.port = "80"
	this.timeout = time.Second
}

func (this *SuiteTCPConnect) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPConnect) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPConnect) TestNewTCPConnect() {
	assert.NotNil(this.T(), NewTCPConnect(this.ip, this.port, this.timeout))
}

func (this *SuiteTCPConnect) TestStart() {
	valid := false
	target := NewTCPConnect(this.ip, this.port, this.timeout)
	target.Start(func(session Sessioner, err error) {
		if session != nil && err == nil {
			valid = true
			fmt.Printf("remote addr: %s\n", session.RemoteAddr().String())
			fmt.Printf("local addr: %s\n", session.LocalAddr().String())
		} // if
	})
	assert.True(this.T(), valid)

	valid = false
	target = NewTCPConnect("!?", this.port, this.timeout)
	target.Start(func(session Sessioner, err error) {
		valid = session != nil && err == nil
	})
	assert.False(this.T(), valid)

	valid = false
	target = NewTCPConnect(this.ip, "3000", this.timeout) // 故意連線到不開放的埠號才會引發錯誤
	target.Start(func(session Sessioner, err error) {
		valid = session != nil && err == nil
	})
	assert.False(this.T(), valid)
}

func (this *SuiteTCPConnect) TestAddress() {
	target := NewTCPConnect(this.ip, this.port, this.timeout)
	assert.Equal(this.T(), net.JoinHostPort(this.ip, this.port), target.Address())
}
