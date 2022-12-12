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

func TestNetmgr(t *testing.T) {
	suite.Run(t, new(SuiteNetmgr))
}

type SuiteNetmgr struct {
	suite.Suite
	testdata.TestEnv
	ip      string
	port    string
	timeout time.Duration
}

func (this *SuiteNetmgr) SetupSuite() {
	this.Change("test-nets-netmgr")
	this.ip = ""
	this.port = "3000"
	this.timeout = time.Second
}

func (this *SuiteNetmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteNetmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteNetmgr) TestNewNetmgr() {
	assert.NotNil(this.T(), NewNetmgr(nil))
}

func (this *SuiteNetmgr) TestAddConnect() {
	tester := newNetmgrTester()
	target := NewNetmgr(tester.Error)
	target.AddConnect(NewTCPConnect("google.com", "80", this.timeout), tester.Prepare)
	assert.True(this.T(), tester.wait())
	assert.True(this.T(), tester.valid())
	target.Stop()
	tester.get().StopWait()
}

func (this *SuiteNetmgr) TestAddListen() {
	testerl := newNetmgrTester()
	target := NewNetmgr(testerl.Error)
	target.AddListen(NewTCPListen(this.ip, this.port), testerl.Prepare)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Connect(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())
	target.Stop()
	testerl.get().StopWait()
	testerc.get().StopWait()
}

func (this *SuiteNetmgr) TestNetmgr() {
	testerl := newNetmgrTester()
	target := NewNetmgr(testerl.Error)
	target.AddListen(NewTCPListen(this.ip, this.port), testerl.Prepare)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Connect(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	time.Sleep(this.timeout)
	sessionID := testerl.get().SessionID()
	assert.Equal(this.T(), testerl.get(), target.GetSession(sessionID))
	status := target.Status()
	assert.Len(this.T(), status.Listen, 1)
	assert.Equal(this.T(), 1, status.session)
	target.StopSession(sessionID)
	assert.Nil(this.T(), target.GetSession(sessionID))

	target.Stop()
	testerl.get().StopWait()
	testerc.get().StopWait()
}

func (this *SuiteNetmgr) TestComplete() {
	validPrepare := false
	validFailure := false
	target := newComplete(
		"127.0.0.1",
		func(_ Sessioner) (coder Coder, reactor Reactor) {
			validPrepare = true
			return nil, nil
		},
		func(err error) {
			validFailure = true
		},
		newSessionmgr(),
	)
	assert.NotNil(this.T(), this)
	target.complete(nil, fmt.Errorf("error"))
	assert.True(this.T(), validFailure)
	target.complete(&emptySession{}, nil)
	assert.True(this.T(), validPrepare)
}

func (this *SuiteNetmgr) TestListenmgr() {
	target := newListenmgr()
	assert.NotNil(this.T(), target)
	target.add(NewTCPListen("127.0.0.1", "1"))
	target.add(NewTCPListen("127.0.0.2", "2"))
	target.add(NewTCPListen("127.0.0.3", "3"))
	assert.ElementsMatch(this.T(), []string{"127.0.0.1:1", "127.0.0.2:2", "127.0.0.3:3"}, target.address())
	target.clear()
}

func (this *SuiteNetmgr) TestSessionmgr() {
	target := newSessionmgr()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), SessionID(1), target.add(&emptySession{}))
	assert.Equal(this.T(), SessionID(2), target.add(&emptySession{}))
	assert.Equal(this.T(), SessionID(3), target.add(&emptySession{}))
	assert.Equal(this.T(), 3, target.count())
	assert.NotNil(this.T(), target.get(1))
	target.del(1)
	assert.Equal(this.T(), 2, target.count())
	assert.Nil(this.T(), target.get(1))
	target.clear()
	assert.Equal(this.T(), 0, target.count())
	assert.Nil(this.T(), target.get(2))
	assert.Nil(this.T(), target.get(3))
}
