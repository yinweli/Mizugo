package helps

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestSync(t *testing.T) {
	suite.Run(t, new(SuiteSync))
}

type SuiteSync struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteSync) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-sync"))
}

func (this *SuiteSync) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteSync) TestSyncOnce() {
	target := SyncOnce{}
	valid := atomic.Int64{}
	validFunc := func() {
		valid.Add(1)
	}
	go target.Do(validFunc)
	go target.Do(validFunc)
	trials.WaitTimeout()
	assert.Equal(this.T(), int64(1), valid.Load())
	assert.True(this.T(), target.Done())
}

func (this *SuiteSync) TestSyncAttr() {
	data := "data"
	target := SyncAttr[string]{}
	target.Set(data)
	assert.Equal(this.T(), data, target.Get())
}
