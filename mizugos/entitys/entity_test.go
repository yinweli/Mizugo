package entitys

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestEntity(t *testing.T) {
	suite.Run(t, new(SuiteEntity))
}

type SuiteEntity struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEntity) SetupSuite() {
	this.Change("test-entitys-entity")
}

func (this *SuiteEntity) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntity) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntity) TestNewEntity() {
	assert.NotNil(this.T(), NewEntity(EntityID(1)))
}

func (this *SuiteEntity) TestInitialize() {
	entityID := EntityID(1)
	target := NewEntity(entityID)
	module := newModuleTester(ModuleID(1))

	assert.Nil(this.T(), target.AddModule(module))
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Initialize())
	assert.Equal(this.T(), int64(1), module.awake.Load())
	assert.Equal(this.T(), int64(1), module.start.Load())

	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行
}

func (this *SuiteEntity) TestEntity() {
	entityID := EntityID(1)
	target := NewEntity(entityID)

	assert.Nil(this.T(), target.SetProcess(procs.NewSimple()))
	assert.Nil(this.T(), target.Initialize())

	bundle := target.Bundle()

	assert.NotNil(this.T(), bundle.Encode)
	assert.NotNil(this.T(), bundle.Decode)
	assert.NotNil(this.T(), bundle.Receive)
	assert.NotNil(this.T(), bundle.AfterSend)
	assert.NotNil(this.T(), bundle.AfterRecv)
	assert.Equal(this.T(), entityID, target.EntityID())
	assert.True(this.T(), target.Enable())

	target.Finalize()
}

func (this *SuiteEntity) TestSession() {
	target := NewEntity(EntityID(1))
	conn, _ := net.Dial("tcp", net.JoinHostPort("google.com", "80"))
	session := nets.NewTCPSession(conn)

	assert.Nil(this.T(), target.SetSession(session))
	assert.Equal(this.T(), session, target.GetSession())

	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.SetSession(session))

	target.Send("message")
	assert.Equal(this.T(), session.RemoteAddr(), target.RemoteAddr())
	assert.Equal(this.T(), session.LocalAddr(), target.LocalAddr())

	target.Finalize()
}

func (this *SuiteEntity) TestProcess() {
	target := NewEntity(EntityID(1))
	process := procs.NewSimple()

	assert.Nil(this.T(), target.SetProcess(process))
	assert.Equal(this.T(), process, target.GetProcess())
	target.AddMessage(procs.MessageID(1), func(message any) {
		// do nothing
	})
	target.DelMessage(procs.MessageID(1))

	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.SetProcess(process))

	target.Finalize()
}

func (this *SuiteEntity) TestModule() {
	target := NewEntity(EntityID(1))
	module1 := newModuleTester(ModuleID(1))
	module2 := newModuleTester(ModuleID(2))

	assert.Nil(this.T(), target.AddModule(module1))
	assert.NotNil(this.T(), target.GetModule(module1.ModuleID()))
	assert.NotNil(this.T(), target.AddModule(module1))

	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.AddModule(module2))
	target.Finalize()
}

func (this *SuiteEntity) TestEvent() {
	target := NewEntity(EntityID(1))
	nameOnce := "event once"
	valueOnce := "value once"
	validOnce := atomic.Bool{}
	nameFixed := "event fixed"
	valueFixed := "value fixed"
	validFixed := atomic.Bool{}

	eventIDOnce, err := target.SubEvent(nameOnce, func(param any) {
		validOnce.Store(param.(string) == valueOnce)
	})
	assert.Nil(this.T(), err)

	eventIDFixed, err := target.SubEvent(nameFixed, func(param any) {
		validFixed.Store(param.(string) == valueFixed)
	})
	assert.Nil(this.T(), err)

	_, err = target.SubEvent(EventFinalize, nil)
	assert.NotNil(this.T(), err)

	assert.Nil(this.T(), target.Initialize())
	target.PubOnceEvent(nameOnce, valueOnce)
	target.PubFixedEvent(nameFixed, valueFixed, time.Millisecond)

	time.Sleep(testdata.Timeout * 2) // 多等一下讓定時事件發生
	assert.True(this.T(), validOnce.Load())
	assert.True(this.T(), validFixed.Load())

	target.UnsubEvent(eventIDOnce)
	target.UnsubEvent(eventIDFixed)

	target.Finalize()
}
