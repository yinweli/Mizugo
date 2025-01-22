package redmos

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdIncr(t *testing.T) {
	suite.Run(t, new(SuiteCmdIncr))
}

type SuiteCmdIncr struct {
	suite.Suite
	trials.Catalog
	meta  metaIncr
	major *Major
	minor *Minor
}

func (this *SuiteCmdIncr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdincr"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdincr")
}

func (this *SuiteCmdIncr) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdIncr) TestIncr() {
	this.meta.table = true
	majorSubmit := this.major.Submit()
	minorSubmit := this.minor.Submit()
	key := "redis+mongo"
	data := &IncrData{}

	target := &Incr{MinorEnable: true, Meta: &this.meta, Key: key, Incr: 1, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.Nil(this.T(), target.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	assert.Nil(this.T(), target.Complete())
	_ = minorSubmit.Exec(context.Background())
	assert.Equal(this.T(), int64(1), target.Data.Data)
	assert.True(this.T(), trials.MongoCompare[IncrData](this.minor.Database(), this.meta.MinorTable(), MongoKey, this.meta.MinorKey(key), data))

	target = &Incr{MinorEnable: true, Meta: nil, Key: key, Incr: 1, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	target = &Incr{MinorEnable: true, Meta: &this.meta, Key: "", Incr: 1, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())

	this.meta.table = false
	target = &Incr{MinorEnable: true, Meta: &this.meta, Key: key, Incr: 1, Data: data}
	target.Initialize(context.Background(), majorSubmit, minorSubmit)
	assert.NotNil(this.T(), target.Prepare())
}

func (this *SuiteCmdIncr) TestDuplicate() {
	this.meta.table = true
	count := 4
	total := atomic.Int64{}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer waitGroup.Done()
			trials.WaitTimeout()

			for i := 0; i < 100; i++ {
				majorSubmit := this.major.Submit()
				minorSubmit := this.minor.Submit()
				incr := &Incr{MinorEnable: true, Meta: &this.meta, Key: "cmdincr+duplicate", Incr: 1, Data: &IncrData{}}
				incr.Initialize(context.Background(), majorSubmit, minorSubmit)
				_ = incr.Prepare()
				_, _ = majorSubmit.Exec(context.Background())
				_ = incr.Complete()
				total.Add(incr.Data.Data)
			} // for
		}()
	} // for

	waitGroup.Wait()
	assert.Equal(this.T(), int64(80200), total.Load())
}

type metaIncr struct {
	table bool
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
