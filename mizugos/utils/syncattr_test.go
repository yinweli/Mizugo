package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestSyncAttr(t *testing.T) {
	suite.Run(t, new(SuiteSyncAttr))
}

type SuiteSyncAttr struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteSyncAttr) SetupSuite() {
	this.Change("test-utils-syncattr")
}

func (this *SuiteSyncAttr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteSyncAttr) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteSyncAttr) TestSetGet() {
	data := "data"
	target := SyncAttr[string]{}
	target.Set(data)
	assert.Equal(this.T(), data, target.Get())
}
