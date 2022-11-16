package entitys

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

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
	this.Change("test-entity")
}

func (this *SuiteEntity) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntity) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntity) TestNewEntity() {
	assert.NotNil(this.T(), NewEntity(EntityID(1), "entity"))
}

func (this *SuiteEntity) TestEntity() {
	target := NewEntity(EntityID(1), "entity")
	target.initialize()
	defer target.finalize()

	assert.Equal(this.T(), EntityID(1), target.EntityID())
	assert.Equal(this.T(), "entity", target.Name())
}

func (this *SuiteEntity) TestAddModule() {
	target := NewEntity(EntityID(1), "entity")
	target.initialize()
	defer target.finalize()

	assert.Nil(this.T(), target.AddModule(NewModule(ModuleID(1), "module")))
	assert.NotNil(this.T(), target.GetModule(ModuleID(1)))

	assert.NotNil(this.T(), target.AddModule(NewModule(ModuleID(1), "module")))
}

func (this *SuiteEntity) TestDelModule() {
	target := NewEntity(EntityID(1), "entity")
	target.initialize()
	defer target.finalize()

	assert.Nil(this.T(), target.AddModule(NewModule(ModuleID(1), "module")))
	assert.NotNil(this.T(), target.DelModule(ModuleID(1)))
	assert.Nil(this.T(), target.GetModule(ModuleID(1)))

	assert.Nil(this.T(), target.DelModule(ModuleID(1)))
}

func (this *SuiteEntity) TestGetModule() {
	target := NewEntity(EntityID(1), "entity")
	target.initialize()
	defer target.finalize()

	assert.Nil(this.T(), target.AddModule(NewModule(ModuleID(1), "module")))
	assert.NotNil(this.T(), target.GetModule(ModuleID(1)))

	assert.Nil(this.T(), target.GetModule(ModuleID(2)))
}

func (this *SuiteEntity) TestPubOnceEvent() {
	target := NewEntity(EntityID(1), "entity")
	target.initialize()
	defer target.finalize()

	valid := atomic.Bool{}
	target.SubEvent("event", func(param any) {
		if param.(string) == "pubonce" {
			valid.Store(true)
		} // if
	})
	target.PubOnceEvent("event", "pubonce")
	time.Sleep(time.Millisecond * 10)
	assert.True(this.T(), valid.Load())
}

func (this *SuiteEntity) TestPubFixedEvent() {
	target := NewEntity(EntityID(1), "entity")
	target.initialize()
	defer target.finalize()

	valid := atomic.Int64{}
	target.SubEvent("event", func(param any) {
		if param.(string) == "pubfixed" {
			valid.Add(1)
		} // if
	})
	fixed := target.PubFixedEvent("event", "pubfixed", time.Millisecond)
	defer fixed.Stop()
	time.Sleep(time.Millisecond * 100)
	assert.Greater(this.T(), valid.Load(), int64(0))
}

func (this *SuiteEntity) TestInitialize() {
	target := NewEntity(EntityID(1), "entity")
	module := &testModule{
		Module: Module{
			moduleID: ModuleID(1),
			name:     "module",
		},
	}
	assert.Nil(this.T(), target.AddModule(module))
	target.initialize()
	time.Sleep(updateInterval * 2)
	target.finalize()

	assert.True(this.T(), module.awake.Load())
	assert.True(this.T(), module.start.Load())
	assert.True(this.T(), module.dispose.Load())
	assert.True(this.T(), module.update.Load())
}

type testModule struct {
	Module
	awake   atomic.Bool
	start   atomic.Bool
	dispose atomic.Bool
	update  atomic.Bool
}

func (this *testModule) Awake() {
	this.awake.Store(true)
}

func (this *testModule) Start() {
	this.start.Store(true)
}

func (this *testModule) Dispose() {
	this.dispose.Store(true)
}

func (this *testModule) Update() {
	this.update.Store(true)
}
