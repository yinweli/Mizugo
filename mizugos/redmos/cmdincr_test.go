package redmos

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdIncr(t *testing.T) {
	suite.Run(t, new(SuiteCmdIncr))
}

type SuiteCmdIncr struct {
	suite.Suite
	testdata.Env
	meta  metaIncr
	major *Major
	minor *Minor
}

func (this *SuiteCmdIncr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-cmdincr")
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdincr")
}

func (this *SuiteCmdIncr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdIncr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCmdIncr) TestIncr() {
	this.meta.table = true
	this.meta.field = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	data := &dataIncr{Field: "redis+mongo", Value: 1}

	target := &Incr{Meta: &this.meta, MinorEnable: true, Key: data.Field, Data: &IncrData{Incr: 1, Value: 0}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(ctxs.Get().Ctx())
	assert.Equal(this.T(), int64(1), target.Data.Value)
	assert.True(this.T(), testdata.MongoCompare[dataIncr](this.minor.Database(), this.meta.MinorTable(), this.meta.MinorField(), this.meta.MinorKey(data.Field), data))

	target = &Incr{Meta: nil, MinorEnable: true, Key: data.Field, Data: &IncrData{Incr: 1, Value: 0}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Incr{Meta: &this.meta, MinorEnable: true, Key: "", Data: &IncrData{Incr: 1, Value: 0}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Incr{Meta: &this.meta, MinorEnable: true, Key: data.Field, Data: nil}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	this.meta.field = true
	target = &Incr{Meta: &this.meta, MinorEnable: true, Key: data.Field, Data: &IncrData{Incr: 1, Value: 0}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = true
	this.meta.field = false
	target = &Incr{Meta: &this.meta, MinorEnable: true, Key: data.Field, Data: &IncrData{Incr: 1, Value: 0}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Incr{Meta: nil, MinorEnable: true, Key: data.Field, Data: &IncrData{Incr: 1, Value: 0}}
	target.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Complete())
}

func (this *SuiteCmdIncr) TestDuplicate() {
	this.meta.table = true
	this.meta.field = true
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(4)
	count := atomic.Int64{}
	check := func() {
		testdata.WaitTimeout()

		for i := 0; i < 250; i++ {
			majorSubmit := this.major.Submit()
			minorSubmit := this.minor.Submit()
			data := &dataIncr{Field: "duplicate", Value: 1}
			incr := &Incr{Meta: &this.meta, MinorEnable: false, Key: data.Field, Data: &IncrData{Incr: 1, Value: 0}}
			incr.Initialize(ctxs.Get().Ctx(), majorSubmit, minorSubmit)
			_ = incr.Prepare()
			_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
			_ = incr.Complete()
			_ = minorSubmit.Exec(ctxs.Get().Ctx())
			count.Add(incr.Data.Value)
		} // for

		waitGroup.Done()
	}

	go check()
	go check()
	go check()
	go check()
	waitGroup.Wait()
	assert.Equal(this.T(), int64(500500), count.Load())
}

type metaIncr struct {
	table bool
	field bool
}

func (this *metaIncr) MajorKey(key any) string {
	return fmt.Sprintf("cmdincr:%v", key)
}

func (this *metaIncr) MinorKey(key any) string {
	return fmt.Sprintf("%v", key)
}

func (this *metaIncr) MinorTable() string {
	if this.table {
		return "cmdincr"
	} // if

	return ""
}

func (this *metaIncr) MinorField() string {
	if this.field {
		return "field"
	} // if

	return ""
}

type dataIncr struct {
	Field string `bson:"field"`
	Value int64  `bson:"value"`
}
