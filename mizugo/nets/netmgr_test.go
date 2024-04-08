package nets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestNetmgr(t *testing.T) {
	suite.Run(t, new(SuiteNetmgr))
}

type SuiteNetmgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteNetmgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-nets-netmgr")
}

func (this *SuiteNetmgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteNetmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteNetmgr) TestAddConnectTCP() {
	addr := host{ip: "google.com", port: "80"}
	test := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	connectID := target.AddConnectTCP(addr.ip, addr.port, testdata.Timeout, test.bind, test.unbind, test.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), test.valid())
	assert.NotNil(this.T(), target.GetConnect(connectID))
	target.DelConnect(connectID)
	assert.Nil(this.T(), target.GetConnect(connectID))

	time.Sleep(testdata.Timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestAddListenTCP() {
	addr := host{port: "9000"}
	test := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	listenID := target.AddListenTCP(addr.ip, addr.port, test.bind, test.unbind, test.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), test.valid())
	assert.NotNil(this.T(), target.GetListen(listenID))
	target.DelListen(listenID)
	assert.Nil(this.T(), target.GetListen(listenID))

	time.Sleep(testdata.Timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestStop() {
	addr := host{ip: "google.com", port: "80"}
	test := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	target.AddConnectTCP(addr.ip, addr.port, testdata.Timeout, test.bind, test.unbind, test.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), test.validSession())

	time.Sleep(testdata.Timeout)
	target.Stop()

	time.Sleep(testdata.Timeout)
	assert.False(this.T(), test.validSession())
}

func (this *SuiteNetmgr) TestStatus() {
	addr := host{port: "9000"}
	testl := newTester(true, true, true)
	testc := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	target.AddListenTCP(addr.ip, addr.port, testl.bind, testl.unbind, testl.wrong)
	target.AddConnectTCP(addr.ip, addr.port, testdata.Timeout, testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())
	status := target.Status()
	assert.Equal(this.T(), 1, status.Connect)
	assert.Equal(this.T(), 1, status.Listen)
	assert.Equal(this.T(), 2, status.Session)

	time.Sleep(testdata.Timeout)
	target.Stop()
}

func (this *SuiteNetmgr) TestConnectmgr() {
	connect1 := &emptyConnect{value: 1}
	connect2 := &emptyConnect{value: 2}
	connect3 := &emptyConnect{value: 3}
	target := newConnectmgr()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), ConnectID(1), target.add(connect1))
	assert.Equal(this.T(), ConnectID(2), target.add(connect2))
	assert.Equal(this.T(), ConnectID(3), target.add(connect3))
	assert.Equal(this.T(), connect1, target.get(ConnectID(1)))
	assert.Equal(this.T(), connect2, target.get(ConnectID(2)))
	assert.Equal(this.T(), connect3, target.get(ConnectID(3)))
	assert.Equal(this.T(), 3, target.count())
	target.del(ConnectID(1))
	assert.Equal(this.T(), nil, target.get(ConnectID(1)))
	assert.Equal(this.T(), 2, target.count())
	target.clear()
	assert.Equal(this.T(), nil, target.get(ConnectID(2)))
	assert.Equal(this.T(), nil, target.get(ConnectID(3)))
	assert.Equal(this.T(), 0, target.count())
}

func (this *SuiteNetmgr) TestListenmgr() {
	listen1 := &emptyListen{value: 1}
	listen2 := &emptyListen{value: 2}
	listen3 := &emptyListen{value: 3}
	target := newListenmgr()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), ListenID(1), target.add(listen1))
	assert.Equal(this.T(), ListenID(2), target.add(listen2))
	assert.Equal(this.T(), ListenID(3), target.add(listen3))
	assert.Equal(this.T(), listen1, target.get(ListenID(1)))
	assert.Equal(this.T(), listen2, target.get(ListenID(2)))
	assert.Equal(this.T(), listen3, target.get(ListenID(3)))
	assert.Equal(this.T(), 3, target.count())
	target.del(ListenID(1))
	assert.Equal(this.T(), nil, target.get(ListenID(1)))
	assert.Equal(this.T(), 2, target.count())
	target.clear()
	assert.Equal(this.T(), nil, target.get(ListenID(2)))
	assert.Equal(this.T(), nil, target.get(ListenID(3)))
	assert.Equal(this.T(), 0, target.count())
}

func (this *SuiteNetmgr) TestSessionmgr() {
	session1 := &emptySession{value: 1}
	session2 := &emptySession{value: 2}
	session3 := &emptySession{value: 3}
	target := newSessionmgr()
	assert.NotNil(this.T(), target)
	target.add(session1)
	target.add(session2)
	target.add(session3)
	assert.Equal(this.T(), 3, target.count())
	target.del(session3)
	assert.Equal(this.T(), 2, target.count())
	target.clear()
	assert.Equal(this.T(), 0, target.count())
}
