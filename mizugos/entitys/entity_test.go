package entitys

import (
	"net"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/nets"
	"github.com/yinweli/Mizugo/v2/mizugos/procs"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestEntity(t *testing.T) {
	suite.Run(t, new(SuiteEntity))
}

type SuiteEntity struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteEntity) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-entitys-entity"))
}

func (this *SuiteEntity) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteEntity) TestEntity() {
	target := NewEntity(EntityID(1))
	this.NotNil(target)
	this.Nil(target.Initialize(nil))
	this.Equal(EntityID(1), target.EntityID())
	this.True(target.Enable())
	target.Finalize()
}

func (this *SuiteEntity) TestInitialize() {
	target := NewEntity(EntityID(1))
	target.Finalize() // 初始化前執行, 這次應該不執行
	this.Nil(target.Initialize(nil))
	this.NotNil(target.Initialize(nil)) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行

	target = NewEntity(EntityID(1))
	module := newTestModule(true, true, ModuleID(1))
	this.Nil(target.AddModule(module))
	this.Nil(target.Initialize(nil))
	this.Equal(int64(1), module.awakeCount.Load())
	this.Equal(int64(1), module.startCount.Load())
	target.Finalize()

	target = NewEntity(EntityID(1))
	this.Nil(target.AddModule(newTestModule(false, true, ModuleID(1))))
	this.NotNil(target.Initialize(nil))
	target.Finalize()

	target = NewEntity(EntityID(1))
	this.Nil(target.AddModule(newTestModule(true, false, ModuleID(1))))
	this.NotNil(target.Initialize(nil))
	target.Finalize()
}

func (this *SuiteEntity) TestModule() {
	target := NewEntity(EntityID(1))
	module1 := newTestModule(true, true, ModuleID(1))
	module2 := newTestModule(true, true, ModuleID(2))
	this.Nil(target.AddModule(module1))
	this.NotNil(target.GetModule(module1.ModuleID()))
	this.Nil(target.Initialize(nil))
	this.NotNil(target.AddModule(module2))
	target.Finalize()
}

func (this *SuiteEntity) TestEvent() {
	target := NewEntity(EntityID(1))
	this.Nil(target.Initialize(nil))

	onceValue := "once"
	onceValid := atomic.Bool{}
	onceSubID := target.Subscribe(onceValue, func(param any) {
		onceValid.Store(param.(string) == onceValue)
	})
	this.NotNil(onceSubID)

	delayValue := "delay"
	delayValid := atomic.Bool{}
	delaySubID := target.Subscribe(delayValue, func(param any) {
		delayValid.Store(param.(string) == delayValue)
	})
	this.NotNil(delaySubID)

	fixedValue := "fixed"
	fixedValid := atomic.Bool{}
	fixedSubID := target.Subscribe(fixedValue, func(param any) {
		fixedValid.Store(param.(string) == fixedValue)
	})
	this.NotNil(fixedSubID)

	target.PublishOnce(onceValue, onceValue)
	target.PublishDelay(onceValue, onceValue, trials.Timeout)
	target.PublishFixed(fixedValue, fixedValue, trials.Timeout)

	trials.WaitTimeout(trials.Timeout * 2) // 多等一下讓定時事件發生
	this.True(onceValid.Load())
	this.True(fixedValid.Load())
	target.Unsubscribe(onceSubID)
	target.Unsubscribe(fixedSubID)
	target.Finalize()
}

func (this *SuiteEntity) TestProcess() {
	target := NewEntity(EntityID(1))
	process := procs.NewJson()
	this.Nil(target.SetProcess(process))
	this.Nil(target.Initialize(nil))
	this.NotNil(target.SetProcess(process))
	target.AddMessage(1, func(_ any) {})
	this.NotNil(target.GetMessage(1))
	target.DelMessage(1)
	this.Nil(target.GetMessage(1))
	target.Finalize()
}

func (this *SuiteEntity) TestSession() {
	target := NewEntity(EntityID(1))
	conn, _ := net.Dial("tcp", net.JoinHostPort("google.com", "80"))
	session := nets.NewTCPSession(conn)
	this.Nil(target.SetSession(session))
	this.Nil(target.Initialize(nil))
	this.NotNil(target.SetSession(session))
	target.Send("message")
	this.Equal(session.RemoteAddr(), target.RemoteAddr())
	this.Equal(session.LocalAddr(), target.LocalAddr())
	target.Stop()
	target.StopWait()
	target.Finalize()
}
