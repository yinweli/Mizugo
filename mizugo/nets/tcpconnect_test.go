package nets

import (
	"fmt"
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
	port    int
	timeout time.Duration
}

func (this *SuiteTCPConnect) SetupSuite() {
	this.Change("test-tcpconnect")
	this.ip = "google.com"
	this.port = 80
	this.timeout = time.Second * 5
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
	// TODO: 單元測試
}

func (this *SuiteTCPConnect) TestAddress() {
	target := NewTCPConnect(this.ip, this.port, this.timeout)
	addr, err := target.Address()
	assert.Nil(this.T(), err)
	assert.NotEmpty(this.T(), addr.String())
	fmt.Printf("%s:%d >>> %s\n", this.ip, this.port, addr.String())
}
