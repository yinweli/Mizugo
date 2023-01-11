package entitys

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/mizugos/events"
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
	entityID EntityID
}

func (this *SuiteEntity) SetupSuite() {
	this.Change("test-entitys-entity")
	this.entityID = EntityID(1)
}

func (this *SuiteEntity) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntity) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntity) TestNewEntity() {
	assert.NotNil(this.T(), NewEntity(this.entityID))
}

func (this *SuiteEntity) TestInitialize() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Initialize()) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))
	module := newModuleTester(true, true, ModuleID(1))
	assert.Nil(this.T(), target.AddModule(module))
	assert.Nil(this.T(), target.Initialize())
	assert.Equal(this.T(), int64(1), module.awakeCount.Load())
	assert.Equal(this.T(), int64(1), module.startCount.Load())
	target.Finalize()

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))
	assert.Nil(this.T(), target.AddModule(newModuleTester(false, true, ModuleID(1))))
	assert.NotNil(this.T(), target.Initialize())

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))
	assert.Nil(this.T(), target.AddModule(newModuleTester(true, false, ModuleID(1))))
	assert.NotNil(this.T(), target.Initialize())
}

func (this *SuiteEntity) TestEntity() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))
	assert.Nil(this.T(), target.SetProcess(procs.NewSimple()))
	assert.Nil(this.T(), target.Initialize())

	bundle := target.Bundle()
	assert.NotNil(this.T(), bundle.Encode)
	assert.NotNil(this.T(), bundle.Decode)
	assert.NotNil(this.T(), bundle.Receive)
	assert.NotNil(this.T(), bundle.AfterSend)
	assert.NotNil(this.T(), bundle.AfterRecv)
	assert.Equal(this.T(), this.entityID, target.EntityID())
	assert.True(this.T(), target.Enable())

	target.Finalize()
}

func (this *SuiteEntity) TestModule() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))

	assert.NotNil(this.T(), target.GetModulemgr())

	module1 := newModuleTester(true, true, ModuleID(1))
	module2 := newModuleTester(true, true, ModuleID(2))
	assert.Nil(this.T(), target.AddModule(module1))
	assert.NotNil(this.T(), target.GetModule(module1.ModuleID()))
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.AddModule(module2))

	target.Finalize()
}

func (this *SuiteEntity) TestEvent() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))
	assert.Nil(this.T(), target.Initialize())

	assert.NotNil(this.T(), target.GetEventmgr())

	onceValue := "once"
	onceValid := atomic.Bool{}
	onceSubID, err := target.Subscribe(onceValue, func(param any) {
		onceValid.Store(param.(string) == onceValue)
	})
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), onceSubID)

	fixedValue := "fixed"
	fixedValid := atomic.Bool{}
	fixedSubID, err := target.Subscribe(fixedValue, func(param any) {
		fixedValid.Store(param.(string) == fixedValue)
	})
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), fixedSubID)

	_, err = target.Subscribe(EventFinalize, nil)
	assert.NotNil(this.T(), err)

	target.PublishOnce(onceValue, onceValue)
	target.PublishFixed(fixedValue, fixedValue, time.Millisecond)

	time.Sleep(testdata.Timeout * 2) // 多等一下讓定時事件發生
	assert.True(this.T(), onceValid.Load())
	assert.True(this.T(), fixedValid.Load())

	target.Unsubscribe(onceSubID)
	target.Unsubscribe(fixedSubID)
	target.Finalize()
}

func (this *SuiteEntity) TestProcess() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))

	process := procs.NewSimple()
	assert.Nil(this.T(), target.SetProcess(process))
	assert.Equal(this.T(), process, target.GetProcess())
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.SetProcess(process))

	target.AddMessage(procs.MessageID(1), func(message any) {})
	target.DelMessage(procs.MessageID(1))

	target.Finalize()
}

func (this *SuiteEntity) TestSession() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(1)))

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
