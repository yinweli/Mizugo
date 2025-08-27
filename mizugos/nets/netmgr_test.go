package nets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestNetmgr(t *testing.T) {
	suite.Run(t, new(SuiteNetmgr))
}

type SuiteNetmgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteNetmgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-nets-netmgr"))
}

func (this *SuiteNetmgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteNetmgr) TestAddConnectTCP() {
	addr := host{ip: "google.com", port: "80"}
	test := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	connectID := target.AddConnectTCP(addr.ip, addr.port, trials.Timeout, test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), test.Valid())
	assert.NotNil(this.T(), target.GetConnect(connectID))
	target.DelConnect(connectID)
	assert.Nil(this.T(), target.GetConnect(connectID))

	trials.WaitTimeout()
	target.Stop()
}

func (this *SuiteNetmgr) TestAddListenTCP() {
	addr := host{port: "9000"}
	test := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	listenID := target.AddListenTCP(addr.ip, addr.port, test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), test.Valid())
	assert.NotNil(this.T(), target.GetListen(listenID))
	target.DelListen(listenID)
	assert.Nil(this.T(), target.GetListen(listenID))

	trials.WaitTimeout()
	target.Stop()
}

func (this *SuiteNetmgr) TestStop() {
	addr := host{ip: "google.com", port: "80"}
	test := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	target.AddConnectTCP(addr.ip, addr.port, trials.Timeout, test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), test.ValidSession())

	trials.WaitTimeout()
	target.Stop()

	trials.WaitTimeout()
	assert.False(this.T(), test.ValidSession())
}

func (this *SuiteNetmgr) TestStatus() {
	addr := host{port: "9000"}
	testl := newTester(true, true, true)
	testc := newTester(true, true, true)
	target := NewNetmgr()
	assert.NotNil(this.T(), target)
	target.AddListenTCP(addr.ip, addr.port, testl.Bind, testl.Unbind, testl.Wrong)
	target.AddConnectTCP(addr.ip, addr.port, trials.Timeout, testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.Valid())
	assert.True(this.T(), testc.Valid())
	status := target.Status()
	assert.Equal(this.T(), 1, status.Connect)
	assert.Equal(this.T(), 1, status.Listen)
	assert.Equal(this.T(), 2, status.Session)

	trials.WaitTimeout()
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
