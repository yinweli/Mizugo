package helps

import (
	"sync/atomic"
	"testing"

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
	valid := atomic.Int32{}
	validFunc := func() {
		valid.Add(1)
	}
	target := SyncOnce{}
	go target.Do(validFunc)
	go target.Do(validFunc)
	trials.WaitTimeout()
	this.Equal(int32(1), valid.Load())
	this.True(target.Done())
}

func (this *SuiteSync) TestSyncAttr() {
	target := SyncAttr[string]{}
	target.Set(testdata.Unknown)
	this.Equal(testdata.Unknown, target.Get())
}
