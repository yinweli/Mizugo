package utils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestSync(t *testing.T) {
	suite.Run(t, new(SuiteSync))
}

type SuiteSync struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteSync) SetupSuite() {
	this.Change("test-utils-sync")
}

func (this *SuiteSync) TearDownSuite() {
	this.Restore()
}

func (this *SuiteSync) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteSync) TestSyncOnce() {
	target := SyncOnce{}
	valid := atomic.Int64{}
	validFunc := func() {
		valid.Add(1)
	}

	assert.False(this.T(), target.Done())
	go target.Do(validFunc)
	go target.Do(validFunc)
	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), int64(1), valid.Load())
	assert.True(this.T(), target.Done())
}

func (this *SuiteSync) TestSyncAttr() {
	data := "data"
	target := SyncAttr[string]{}
	target.Set(data)
	assert.Equal(this.T(), data, target.Get())
}
