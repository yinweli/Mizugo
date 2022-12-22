package entitys

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/cores/msgs"
	"github.com/yinweli/Mizugo/cores/nets"
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
	assert.NotNil(this.T(), newEntity(EntityID(1)))
}

func (this *SuiteEntity) TestInitialize() {
	target := newEntity(EntityID(1))
	module := newModuleTester(ModuleID(1))

	assert.Nil(this.T(), target.AddModule(module))
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Initialize())
	assert.True(this.T(), target.Enable())

	time.Sleep(updateInterval * 2) // 為了讓update會被執行, 需要長一點的時間
	assert.Nil(this.T(), target.Finalize())
	assert.NotNil(this.T(), target.Finalize())
	assert.False(this.T(), target.Enable())

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), module.awake.Load())
	assert.True(this.T(), module.start.Load())
	assert.True(this.T(), module.dispose.Load())
	assert.True(this.T(), module.update.Load())
}

func (this *SuiteEntity) TestModule() {
	target := newEntity(EntityID(1))
	module1 := newModuleTester(ModuleID(1))
	module2 := newModuleTester(ModuleID(2))

	assert.Nil(this.T(), target.AddModule(module1))
	assert.NotNil(this.T(), target.GetModule(module1.ModuleID()))
	assert.NotNil(this.T(), target.AddModule(module1))
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.AddModule(module2))
	assert.Nil(this.T(), target.Finalize())
}

func (this *SuiteEntity) TestEvent() {
	target := newEntity(EntityID(1))
	eventOnce := "eventOnce"
	paramOnce := "paramOnce"
	validOnce := atomic.Bool{}
	eventFixed := "eventFixed"
	paramFixed := "paramFixed"
	validFixed := atomic.Int64{}

	assert.Nil(this.T(), target.SubEvent(eventOnce, func(param any) {
		validOnce.Store(param.(string) == paramOnce)
	}))
	assert.Nil(this.T(), target.SubEvent(eventFixed, func(param any) {
		if param.(string) == paramFixed {
			validFixed.Add(1)
		} // if
	}))
	assert.Nil(this.T(), target.Initialize())

	target.PubOnceEvent(eventOnce, paramOnce)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), validOnce.Load())

	fixed := target.PubFixedEvent(eventFixed, paramFixed, time.Millisecond)
	time.Sleep(testdata.Timeout)
	assert.Greater(this.T(), validFixed.Load(), int64(0))
	fixed.Stop()

	assert.NotNil(this.T(), target.SubEvent(eventOnce, func(param any) {
		// do nothing
	}))

	assert.Nil(this.T(), target.Finalize())
}

func (this *SuiteEntity) TestSession() {
	target := newEntity(EntityID(1))
	session := nets.NewTCPSession(nil)

	assert.Nil(this.T(), target.SetSession(session))
	assert.Equal(this.T(), session, target.GetSession())

	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.SetSession(session))
	assert.Nil(this.T(), target.Finalize())
}

func (this *SuiteEntity) TestProcess() {
	target := newEntity(EntityID(1))
	process := msgs.NewStringProc()

	assert.Nil(this.T(), target.SetProcess(process))
	assert.Equal(this.T(), process, target.GetProcess())
	target.AddMsgProc(msgs.MessageID(1), func(messageID msgs.MessageID, message any) {
		// do nothing
	})
	target.DelMsgProc(msgs.MessageID(1))

	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.SetProcess(process))
	assert.Nil(this.T(), target.Finalize())
}
