package nets

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/mizugo/utils"
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
	timeout := utils.NewWaitTimeout(this.timeout)
	target := NewTCPListen(this.ip, this.port)
	go target.Start(func(session Sessioner, err error) {
		if session != nil && err == nil {
			timeout.Done()
			fmt.Printf("remote addr: %s\n", session.RemoteAddr().String())
			fmt.Printf("local addr: %s\n", session.LocalAddr().String())
		} // if
	})
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Start(func(session Sessioner, err error) {})
	assert.True(this.T(), timeout.Wait())
	assert.Nil(this.T(), target.Stop())

	valid := false
	target = NewTCPListen("!?", this.port)
	target.Start(func(session Sessioner, err error) {
		valid = session != nil && err == nil
	})
	assert.False(this.T(), valid)

	valid = false
	target = NewTCPListen("192.168.0.1", this.port) // 故意要監聽錯誤位址才會引發錯誤
	target.Start(func(session Sessioner, err error) {
		valid = session != nil && err == nil
	})
	assert.False(this.T(), valid)
}

func (this *SuiteTCPListen) TestAddress() {
	target := NewTCPListen(this.ip, this.port)
	assert.Equal(this.T(), net.JoinHostPort(this.ip, this.port), target.Address())
}
