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
	assert.Equal(this.T(), entityID, target.EntityID())
	assert.True(this.T(), target.Enable())

	time.Sleep(updateInterval * 2) // 為了讓update會被執行, 需要長一點的時間
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行
	assert.False(this.T(), target.Enable())

	time.Sleep(testdata.Timeout * 2) // 為了給finalize執行, 需要長一點的時間
	assert.Equal(this.T(), int64(1), module.awake.Load())
	assert.Equal(this.T(), int64(1), module.start.Load())
	assert.Greater(this.T(), module.update.Load(), int64(0))
	assert.Equal(this.T(), int64(1), module.dispose.Load())
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
	assert.NotNil(this.T(), target.SubEvent(eventOnce, func(param any) {
		// do nothing
	}))

	target.PubOnceEvent(eventOnce, paramOnce)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), validOnce.Load())

	target.PubFixedEvent(eventFixed, paramFixed, time.Millisecond)
	time.Sleep(testdata.Timeout * 5) // 多等一下讓定時事件發生
	assert.Greater(this.T(), validFixed.Load(), int64(0))

	target.Finalize()
}
