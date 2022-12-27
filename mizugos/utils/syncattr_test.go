package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestSyncAttr(t *testing.T) {
	suite.Run(t, new(SuiteSyncAttr))
}

type SuiteSyncAttr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteSyncAttr) SetupSuite() {
	this.Change("test-utils-syncattr")
}

func (this *SuiteSyncAttr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteSyncAttr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteSyncAttr) TestSetGet() {
	data := "data"
	target := SyncAttr[string]{}
	target.Set(data)
	assert.Equal(this.T(), data, target.Get())
}
