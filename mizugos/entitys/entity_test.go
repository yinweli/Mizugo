package entitys

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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
	testdata.TestLeak
	entityID EntityID
	capacity int
}

func (this *SuiteEntity) SetupSuite() {
	this.Change("test-entitys-entity")
	this.entityID = EntityID(1)
	this.capacity = 1
}

func (this *SuiteEntity) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntity) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteEntity) TestNewEntity() {
	assert.NotNil(this.T(), NewEntity(this.entityID))
}

func (this *SuiteEntity) TestInitialize() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))
	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.Initialize(nil)) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))
	module := newModuleTester(true, true, ModuleID(1))
	assert.Nil(this.T(), target.AddModule(module))
	assert.Nil(this.T(), target.Initialize(nil))
	assert.Equal(this.T(), int64(1), module.awakeCount.Load())
	assert.Equal(this.T(), int64(1), module.startCount.Load())
	target.Finalize()

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))
	assert.Nil(this.T(), target.AddModule(newModuleTester(false, true, ModuleID(1))))
	assert.NotNil(this.T(), target.Initialize(nil))

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))
	assert.Nil(this.T(), target.AddModule(newModuleTester(true, false, ModuleID(1))))
	assert.NotNil(this.T(), target.Initialize(nil))

	target = NewEntity(this.entityID)
	assert.NotNil(this.T(), target.Initialize(nil))

	target = NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.NotNil(this.T(), target.Initialize(nil))
}

func (this *SuiteEntity) TestEntity() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))
	assert.Nil(this.T(), target.SetProcess(procs.NewJson()))
	assert.Nil(this.T(), target.Initialize(nil))

	bundle := target.Bundle()
	assert.NotNil(this.T(), bundle.Encode)
	assert.NotNil(this.T(), bundle.Decode)
	assert.NotNil(this.T(), bundle.Publish)
	assert.Equal(this.T(), this.entityID, target.EntityID())
	assert.True(this.T(), target.Enable())

	target.Finalize()
}

func (this *SuiteEntity) TestModule() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))

	assert.NotNil(this.T(), target.GetModulemgr())

	module1 := newModuleTester(true, true, ModuleID(1))
	module2 := newModuleTester(true, true, ModuleID(2))
	assert.Nil(this.T(), target.AddModule(module1))
	assert.NotNil(this.T(), target.GetModule(module1.ModuleID()))
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.AddModule(module2))

	target.Finalize()
}

func (this *SuiteEntity) TestEvent() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))
	assert.Nil(this.T(), target.Initialize(nil))

	assert.NotNil(this.T(), target.GetEventmgr())

	onceValue := "once"
	onceValid := atomic.Bool{}
	onceSubID := target.Subscribe(onceValue, func(param any) {
		onceValid.Store(param.(string) == onceValue)
	})
	assert.NotNil(this.T(), onceSubID)

	delayValue := "delay"
	delayValid := atomic.Bool{}
	delaySubID := target.Subscribe(delayValue, func(param any) {
		delayValid.Store(param.(string) == delayValue)
	})
	assert.NotNil(this.T(), delaySubID)

	fixedValue := "fixed"
	fixedValid := atomic.Bool{}
	fixedSubID := target.Subscribe(fixedValue, func(param any) {
		fixedValid.Store(param.(string) == fixedValue)
	})
	assert.NotNil(this.T(), fixedSubID)

	target.PublishOnce(onceValue, onceValue)
	target.PublishDelay(onceValue, onceValue, testdata.Timeout)
	target.PublishFixed(fixedValue, fixedValue, testdata.Timeout)

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
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))

	process := procs.NewJson()
	assert.Nil(this.T(), target.SetProcess(process))
	assert.Equal(this.T(), process, target.GetProcess())
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.SetProcess(process))

	target.AddMessage(procs.MessageID(1), func(_ any) {})
	target.DelMessage(procs.MessageID(1))

	target.Finalize()
}

func (this *SuiteEntity) TestSession() {
	target := NewEntity(this.entityID)
	assert.Nil(this.T(), target.SetModulemgr(NewModulemgr()))
	assert.Nil(this.T(), target.SetEventmgr(events.NewEventmgr(this.capacity)))

	conn, _ := net.Dial("tcp", net.JoinHostPort("google.com", "80"))
	session := nets.NewTCPSession(conn)
	assert.Nil(this.T(), target.SetSession(session))
	assert.Equal(this.T(), session, target.GetSession())
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.SetSession(session))

	target.Send("message")
	assert.Equal(this.T(), session.RemoteAddr(), target.RemoteAddr())
	assert.Equal(this.T(), session.LocalAddr(), target.LocalAddr())

	target.Finalize()
}
