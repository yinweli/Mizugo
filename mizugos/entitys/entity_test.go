package entitys

import (
	"net"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.Initialize(nil))
	assert.Equal(this.T(), EntityID(1), target.EntityID())
	assert.True(this.T(), target.Enable())
	target.Finalize()
}

func (this *SuiteEntity) TestInitialize() {
	target := NewEntity(EntityID(1))
	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.Initialize(nil)) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行

	target = NewEntity(EntityID(1))
	module := newModuleTester(true, true, ModuleID(1))
	assert.Nil(this.T(), target.AddModule(module))
	assert.Nil(this.T(), target.Initialize(nil))
	assert.Equal(this.T(), int64(1), module.awakeCount.Load())
	assert.Equal(this.T(), int64(1), module.startCount.Load())
	target.Finalize()

	target = NewEntity(EntityID(1))
	assert.Nil(this.T(), target.AddModule(newModuleTester(false, true, ModuleID(1))))
	assert.NotNil(this.T(), target.Initialize(nil))
	target.Finalize()

	target = NewEntity(EntityID(1))
	assert.Nil(this.T(), target.AddModule(newModuleTester(true, false, ModuleID(1))))
	assert.NotNil(this.T(), target.Initialize(nil))
	target.Finalize()
}

func (this *SuiteEntity) TestModule() {
	target := NewEntity(EntityID(1))
	module1 := newModuleTester(true, true, ModuleID(1))
	module2 := newModuleTester(true, true, ModuleID(2))
	assert.Nil(this.T(), target.AddModule(module1))
	assert.NotNil(this.T(), target.GetModule(module1.ModuleID()))
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.AddModule(module2))
	target.Finalize()
}

func (this *SuiteEntity) TestEvent() {
	target := NewEntity(EntityID(1))
	assert.Nil(this.T(), target.Initialize(nil))

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
	target.PublishDelay(onceValue, onceValue, trials.Timeout)
	target.PublishFixed(fixedValue, fixedValue, trials.Timeout)

	trials.WaitTimeout(trials.Timeout * 2) // 多等一下讓定時事件發生
	assert.True(this.T(), onceValid.Load())
	assert.True(this.T(), fixedValid.Load())
	target.Unsubscribe(onceSubID)
	target.Unsubscribe(fixedSubID)
	target.Finalize()
}

func (this *SuiteEntity) TestProcess() {
	target := NewEntity(EntityID(1))
	process := procs.NewJson()
	assert.Nil(this.T(), target.SetProcess(process))
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.SetProcess(process))
	target.AddMessage(1, func(_ any) {})
	assert.NotNil(this.T(), target.GetMessage(1))
	target.DelMessage(1)
	assert.Nil(this.T(), target.GetMessage(1))
	target.Finalize()
}

func (this *SuiteEntity) TestSession() {
	target := NewEntity(EntityID(1))
	conn, _ := net.Dial("tcp", net.JoinHostPort("google.com", "80"))
	session := nets.NewTCPSession(conn)
	assert.Nil(this.T(), target.SetSession(session))
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.SetSession(session))
	target.Send("message")
	assert.Equal(this.T(), session.RemoteAddr(), target.RemoteAddr())
	assert.Equal(this.T(), session.LocalAddr(), target.LocalAddr())
	target.Stop()
	target.StopWait()
	target.Finalize()
}
