package triggers

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestTriggermgr(t *testing.T) {
	suite.Run(t, new(SuiteTriggermgr))
}

type SuiteTriggermgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteTriggermgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-triggers-triggermgr"))
}

func (this *SuiteTriggermgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteTriggermgr) TestTriggermgr() {
	name := "test"
	function := func() {}
	target := NewTriggermgr()
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.Add(name, function))
	assert.NotNil(this.T(), target.Add(name, function))
	assert.NotNil(this.T(), target.Add(testdata.Unknown, function))
	assert.NotNil(this.T(), target.Get(name))
	target.Finalize()
	trials.WaitTimeout() // 多等一下讓結束完成
}

func (this *SuiteTriggermgr) TestWatch() {
	name := "test"
	count := atomic.Int64{}
	target := NewTriggermgr()
	_ = target.Add(name, func() {
		count.Add(1)
	})
	client := newRedis()
	target.Watch(client, name)
	trials.WaitTimeout() // 多等一下讓監聽完成
	client.Publish(context.Background(), name, name)
	trials.WaitTimeout() // 多等一下讓信號完成
	assert.Equal(this.T(), int64(1), count.Load())
	target.Finalize()
	trials.WaitTimeout() // 多等一下讓結束完成
}

func (this *SuiteTriggermgr) TestTrigger() {
	signal := sync.WaitGroup{}
	signal.Add(1)
	target := &Trigger{exec: func() {
		signal.Done()
	}}

	target.Lock()
	go func() {
		target.Invoke()
	}()
	target.Unlock()
	signal.Wait()
}

func newRedis() redis.UniversalClient {
	option := &redis.UniversalOptions{}
	option.Addrs = append(option.Addrs, testdata.RedisIP)
	client := redis.NewUniversalClient(option)
	return client
}
