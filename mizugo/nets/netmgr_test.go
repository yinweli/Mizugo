package nets

import (
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
	googleIp   string
	googlePort string
	localIp    string
	localPort  string
	timeout    time.Duration
}

func (this *SuiteNetmgr) SetupSuite() {
	this.Change("test-nets-netmgr")
	this.googleIp = "google.com"
	this.googlePort = "80"
	this.localIp = ""
	this.localPort = "3000"
	this.timeout = time.Second
}

func (this *SuiteNetmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteNetmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteNetmgr) TestNewNetmgr() {
	assert.NotNil(this.T(), NewNetmgr())
}

func (this *SuiteNetmgr) TestAddConnect() {
	tester := newSessionTester(true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.googleIp, this.googlePort, this.timeout), tester)

	time.Sleep(this.timeout)
	assert.True(this.T(), tester.validSession())

	time.Sleep(this.timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestAddListen() {
	testerl := newSessionTester(true, true, true)
	target := NewNetmgr()
	target.AddListen(NewTCPListen(this.localIp, this.localPort), testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.localIp, this.localPort, this.timeout)
	go client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.validSession())
	assert.True(this.T(), testerc.valid())

	time.Sleep(this.timeout)
	target.Stop()
	testerc.get().StopWait()
}

func (this *SuiteNetmgr) TestGetSession() {
	tester := newSessionTester(true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.googleIp, this.googlePort, this.timeout), tester)

	time.Sleep(this.timeout)
	assert.True(this.T(), tester.validSession())
	assert.Equal(this.T(), tester.get(), target.GetSession(tester.get().SessionID()))

	time.Sleep(this.timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestStopSession() {
	tester := newSessionTester(true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.googleIp, this.googlePort, this.timeout), tester)

	time.Sleep(this.timeout)
	assert.True(this.T(), tester.validSession())
	target.StopSession(tester.get().SessionID())

	time.Sleep(this.timeout)
	assert.False(this.T(), tester.validSession())

	time.Sleep(this.timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestStop() {
	tester := newSessionTester(true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.googleIp, this.googlePort, this.timeout), tester)

	time.Sleep(this.timeout)
	assert.True(this.T(), tester.validSession())

	time.Sleep(this.timeout)
	target.Stop()

	time.Sleep(this.timeout)
	assert.False(this.T(), tester.validSession())
}

func (this *SuiteNetmgr) TestStatus() {
	testerl := newSessionTester(true, true, true)
	target := NewNetmgr()
	target.AddListen(NewTCPListen(this.localIp, this.localPort), testerl)

	testerc := newSessionTester(true, true, true)
	target.AddConnect(NewTCPConnect(this.googleIp, this.googlePort, this.timeout), testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerc.validSession())
	status := target.Status()
	assert.Len(this.T(), status.Listen, 1)
	assert.Equal(this.T(), 1, status.Session)

	time.Sleep(this.timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestListenmgr() {
	target := newListenmgr()
	assert.NotNil(this.T(), target)
	target.add(NewTCPListen("127.0.0.1", "1"))
	target.add(NewTCPListen("127.0.0.2", "2"))
	target.add(NewTCPListen("127.0.0.3", "3"))
	assert.ElementsMatch(this.T(), target.address(), []string{"127.0.0.1:1", "127.0.0.2:2", "127.0.0.3:3"})
	target.clear()
	assert.ElementsMatch(this.T(), target.address(), []string{})
}

func (this *SuiteNetmgr) TestSessionmgr() {
	target := newSessionmgr()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), SessionID(1), target.add(&emptySession{}))
	assert.Equal(this.T(), SessionID(2), target.add(&emptySession{}))
	assert.Equal(this.T(), SessionID(3), target.add(&emptySession{}))
	assert.NotNil(this.T(), target.get(SessionID(1)))
	assert.NotNil(this.T(), target.get(SessionID(2)))
	assert.NotNil(this.T(), target.get(SessionID(3)))
	assert.Nil(this.T(), target.get(SessionID(4)))
	assert.Equal(this.T(), 3, target.count())
	target.del(1)
	assert.Nil(this.T(), target.get(SessionID(1)))
	assert.Equal(this.T(), 2, target.count())
	target.clear()
	assert.Nil(this.T(), target.get(SessionID(2)))
	assert.Nil(this.T(), target.get(SessionID(3)))
	assert.Equal(this.T(), 0, target.count())
}
