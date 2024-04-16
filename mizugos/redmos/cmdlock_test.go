package redmos

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestCmdLock(t *testing.T) {
	suite.Run(t, new(SuiteCmdLock))
}

type SuiteCmdLock struct {
	suite.Suite
	trials.Catalog
	major *Major
	minor *Minor
}

func (this *SuiteCmdLock) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-cmdlock"))
	this.major, _ = newMajor(testdata.RedisURI)
	this.minor, _ = newMinor(testdata.MongoURI, "cmdlock")
}

func (this *SuiteCmdLock) TearDownSuite() {
	trials.Restore(this.Catalog)
	this.major.DropDB()
	this.major.stop()
	this.minor.DropDB()
	this.minor.stop()
}

func (this *SuiteCmdLock) TestLock() {
	majorSubmit := this.major.Submit()
	key := "lock"

	lock := &Lock{Key: key, time: testdata.RedisTimeout}
	lock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)
	unlock := &Unlock{Key: key}
	unlock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)

	assert.Nil(this.T(), lock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), lock.Complete())

	assert.Nil(this.T(), lock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.NotNil(this.T(), lock.Complete())

	assert.Nil(this.T(), unlock.Prepare())
	_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
	assert.Nil(this.T(), unlock.Complete())

	lock.Key = ""
	assert.NotNil(this.T(), lock.Prepare())

	unlock.Key = ""
	assert.NotNil(this.T(), unlock.Prepare())
}

func (this *SuiteCmdLock) TestDuplicate() {
	key := "lock+duplicate"
	count := 4
	total := atomic.Int64{}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(count + 1)

	go func() {
		defer waitGroup.Done()
		majorSubmit := this.major.Submit()
		lock := &Lock{Key: key, time: testdata.RedisTimeout}
		lock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)
		_ = lock.Prepare()
		_, _ = majorSubmit.Exec(ctxs.Get().Ctx())

		if lock.Complete() != nil {
			return
		} // if

		trials.WaitTimeout(time.Second)

		unlock := &Unlock{Key: key}
		unlock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)
		_ = unlock.Prepare()
		_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
		_ = unlock.Complete()
	}()

	for i := 0; i < count; i++ {
		go func() {
			defer waitGroup.Done()
			trials.WaitTimeout()

			for i := 0; i < 100; i++ {
				majorSubmit := this.major.Submit()
				lock := &Lock{Key: key, time: testdata.RedisTimeout}
				lock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)
				_ = lock.Prepare()
				_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
				_ = lock.Complete()

				if lock.Complete() != nil {
					continue
				} // if

				total.Add(1)

				unlock := &Unlock{Key: key}
				unlock.Initialize(ctxs.Get().Ctx(), majorSubmit, nil)
				_ = unlock.Prepare()
				_, _ = majorSubmit.Exec(ctxs.Get().Ctx())
				_ = unlock.Complete()
			} // for
		}()
	} // for

	waitGroup.Wait()
	assert.Zero(this.T(), total.Load())
}
