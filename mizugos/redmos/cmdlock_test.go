package redmos

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	lock := &Lock{Key: "lock", Token: testdata.Unknown, ttl: testdata.RedisTimeout}
	lock.Initialize(context.Background(), majorSubmit, nil)
	unlock := &Unlock{Key: "lock", Token: testdata.Unknown}
	unlock.Initialize(context.Background(), majorSubmit, nil)

	this.Nil(lock.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(lock.Complete())

	this.Nil(lock.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.NotNil(lock.Complete())

	this.Nil(unlock.Prepare())
	_, _ = majorSubmit.Exec(context.Background())
	this.Nil(unlock.Complete())

	lock.Key = ""
	this.NotNil(lock.Prepare())

	lock.Key = "lock"
	lock.Token = ""
	this.NotNil(lock.Prepare())

	unlock.Key = ""
	this.NotNil(unlock.Prepare())

	unlock.Key = "lock"
	unlock.Token = ""
	this.NotNil(unlock.Prepare())
}

func (this *SuiteCmdLock) TestDuplicate() {
	key := "lock+duplicate"
	count := 4
	total := atomic.Int64{}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(count + 1)

	go func() {
		defer waitGroup.Done()
		token := time.Now().Format("20060102150405.000000000")
		majorSubmit := this.major.Submit()
		lock := &Lock{Key: key, Token: token, ttl: testdata.RedisTimeout}
		lock.Initialize(context.Background(), majorSubmit, nil)
		_ = lock.Prepare()
		_, _ = majorSubmit.Exec(context.Background())

		if lock.Complete() != nil {
			return
		} // if

		trials.WaitTimeout(time.Second)

		unlock := &Unlock{Key: key, Token: token}
		unlock.Initialize(context.Background(), majorSubmit, nil)
		_ = unlock.Prepare()
		_, _ = majorSubmit.Exec(context.Background())
		_ = unlock.Complete()
	}()

	for i := 0; i < count; i++ {
		go func() {
			defer waitGroup.Done()
			trials.WaitTimeout()

			for i := 0; i < 100; i++ {
				token := time.Now().Format("20060102150405.000000000")
				majorSubmit := this.major.Submit()
				lock := &Lock{Key: key, Token: token, ttl: testdata.RedisTimeout}
				lock.Initialize(context.Background(), majorSubmit, nil)
				_ = lock.Prepare()
				_, _ = majorSubmit.Exec(context.Background())
				_ = lock.Complete()

				if lock.Complete() != nil {
					continue
				} // if

				total.Add(1)

				unlock := &Unlock{Key: key, Token: token}
				unlock.Initialize(context.Background(), majorSubmit, nil)
				_ = unlock.Prepare()
				_, _ = majorSubmit.Exec(context.Background())
				_ = unlock.Complete()
			} // for
		}()
	} // for

	waitGroup.Wait()
	this.Zero(total.Load())
}
