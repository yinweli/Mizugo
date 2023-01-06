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
	hostGoogle host
	hostLocal  host
}

func (this *SuiteNetmgr) SetupSuite() {
	this.Change("test-nets-netmgr")
	this.hostGoogle = host{ip: "google.com", port: "80"}
	this.hostLocal = host{ip: "", port: "3000"}
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
	bind := newBindTester(true, true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.hostGoogle.ip, this.hostGoogle.port, testdata.Timeout), bind)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bind.validSession())

	time.Sleep(testdata.Timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestAddListen() {
	bind := newBindTester(true, true, true, true)
	target := NewNetmgr()
	target.AddListen(NewTCPListen(this.hostLocal.ip, this.hostLocal.port), bind)

	done := newDoneTester()
	client := NewTCPConnect(this.hostLocal.ip, this.hostLocal.port, testdata.Timeout)
	client.Connect(done.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bind.validSession())
	assert.True(this.T(), done.valid())

	time.Sleep(testdata.Timeout)
	target.Stop()
	done.get().StopWait()
}

func (this *SuiteNetmgr) TestDelListen() {
	bind := newBindTester(true, true, true, true)
	target := NewNetmgr()
	listenID := target.AddListen(NewTCPListen(this.hostLocal.ip, this.hostLocal.port), bind)
	assert.Equal(this.T(), ListenID(1), listenID)

	time.Sleep(testdata.Timeout)
	target.DelListen(listenID)
	target.Stop()
}

func (this *SuiteNetmgr) TestDelSession() {
	bind := newBindTester(true, true, true, true)
	target := NewNetmgr()
	target.AddListen(NewTCPListen(this.hostLocal.ip, this.hostLocal.port), bind)

	done := newDoneTester()
	client := NewTCPConnect(this.hostLocal.ip, this.hostLocal.port, testdata.Timeout)
	client.Connect(done.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bind.validSession())
	target.DelSession(bind.get().SessionID())

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), bind.validSession())

	time.Sleep(testdata.Timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestGetSession() {
	bind := newBindTester(true, true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.hostGoogle.ip, this.hostGoogle.port, testdata.Timeout), bind)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bind.validSession())
	assert.Equal(this.T(), bind.get(), target.GetSession(bind.get().SessionID()))

	time.Sleep(testdata.Timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestStop() {
	bind := newBindTester(true, true, true, true)
	target := NewNetmgr()
	target.AddConnect(NewTCPConnect(this.hostGoogle.ip, this.hostGoogle.port, testdata.Timeout), bind)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bind.validSession())

	time.Sleep(testdata.Timeout)
	target.Stop()

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), bind.validSession())
}

func (this *SuiteNetmgr) TestStatus() {
	bindl := newBindTester(true, true, true, true)
	target := NewNetmgr()
	target.AddListen(NewTCPListen(this.hostLocal.ip, this.hostLocal.port), bindl)

	bindc := newBindTester(true, true, true, true)
	target.AddConnect(NewTCPConnect(this.hostGoogle.ip, this.hostGoogle.port, testdata.Timeout), bindc)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bindc.validSession())
	status := target.Status()
	assert.Len(this.T(), status.Listen, 1)
	assert.Equal(this.T(), 1, status.Session)

	time.Sleep(testdata.Timeout)
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
