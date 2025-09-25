package nets

import (
	"testing"

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
	test := newTestNet(true, true, true)
	target := NewNetmgr()
	this.NotNil(target)
	connectID := target.AddConnectTCP(addr.ip, addr.port, trials.Timeout, test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	this.True(test.Valid())
	this.NotNil(target.GetConnect(connectID))
	target.DelConnect(connectID)
	this.Nil(target.GetConnect(connectID))

	trials.WaitTimeout()
	target.Stop()
}

func (this *SuiteNetmgr) TestAddListenTCP() {
	addr := host{port: "9000"}
	test := newTestNet(true, true, true)
	target := NewNetmgr()
	this.NotNil(target)
	listenID := target.AddListenTCP(addr.ip, addr.port, test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	this.True(test.Valid())
	this.NotNil(target.GetListen(listenID))
	target.DelListen(listenID)
	this.Nil(target.GetListen(listenID))

	trials.WaitTimeout()
	target.Stop()
}

func (this *SuiteNetmgr) TestStop() {
	addr := host{ip: "google.com", port: "80"}
	test := newTestNet(true, true, true)
	target := NewNetmgr()
	this.NotNil(target)
	target.AddConnectTCP(addr.ip, addr.port, trials.Timeout, test.Bind, test.Unbind, test.Wrong)

	trials.WaitTimeout()
	this.True(test.ValidSession())

	trials.WaitTimeout()
	target.Stop()

	trials.WaitTimeout()
	this.False(test.ValidSession())
}

func (this *SuiteNetmgr) TestStatus() {
	addr := host{port: "9000"}
	testl := newTestNet(true, true, true)
	testc := newTestNet(true, true, true)
	target := NewNetmgr()
	this.NotNil(target)
	target.AddListenTCP(addr.ip, addr.port, testl.Bind, testl.Unbind, testl.Wrong)
	target.AddConnectTCP(addr.ip, addr.port, trials.Timeout, testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testc.Valid())
	connect, listen, session := target.Status()
	this.Equal(1, connect)
	this.Equal(1, listen)
	this.Equal(2, session)

	trials.WaitTimeout()
	target.Stop()
}

func (this *SuiteNetmgr) TestConnectmgr() {
	connect1 := &emptyConnect{value: 1}
	connect2 := &emptyConnect{value: 2}
	connect3 := &emptyConnect{value: 3}
	target := newConnectmgr()
	this.NotNil(target)
	this.Equal(ConnectID(1), target.add(connect1))
	this.Equal(ConnectID(2), target.add(connect2))
	this.Equal(ConnectID(3), target.add(connect3))
	this.Equal(connect1, target.get(ConnectID(1)))
	this.Equal(connect2, target.get(ConnectID(2)))
	this.Equal(connect3, target.get(ConnectID(3)))
	this.Equal(3, target.count())
	target.del(ConnectID(1))
	this.Equal(nil, target.get(ConnectID(1)))
	this.Equal(2, target.count())
	target.clear()
	this.Equal(nil, target.get(ConnectID(2)))
	this.Equal(nil, target.get(ConnectID(3)))
	this.Equal(0, target.count())
}

func (this *SuiteNetmgr) TestListenmgr() {
	listen1 := &emptyListen{value: 1}
	listen2 := &emptyListen{value: 2}
	listen3 := &emptyListen{value: 3}
	target := newListenmgr()
	this.NotNil(target)
	this.Equal(ListenID(1), target.add(listen1))
	this.Equal(ListenID(2), target.add(listen2))
	this.Equal(ListenID(3), target.add(listen3))
	this.Equal(listen1, target.get(ListenID(1)))
	this.Equal(listen2, target.get(ListenID(2)))
	this.Equal(listen3, target.get(ListenID(3)))
	this.Equal(3, target.count())
	target.del(ListenID(1))
	this.Equal(nil, target.get(ListenID(1)))
	this.Equal(2, target.count())
	target.clear()
	this.Equal(nil, target.get(ListenID(2)))
	this.Equal(nil, target.get(ListenID(3)))
	this.Equal(0, target.count())
}

func (this *SuiteNetmgr) TestSessionmgr() {
	session1 := &emptySession{value: 1}
	session2 := &emptySession{value: 2}
	session3 := &emptySession{value: 3}
	target := newSessionmgr()
	this.NotNil(target)
	target.add(session1)
	target.add(session2)
	target.add(session3)
	this.Equal(3, target.count())
	target.del(session3)
	this.Equal(2, target.count())
	target.clear()
	this.Equal(0, target.count())
}
