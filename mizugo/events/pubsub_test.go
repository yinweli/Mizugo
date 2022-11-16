package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestPubsub(t *testing.T) {
	suite.Run(t, new(SuitePubsub))
}

type SuitePubsub struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuitePubsub) SetupSuite() {
	this.Change("test-pubsub")
}

func (this *SuitePubsub) TearDownSuite() {
	this.Restore()
}

func (this *SuitePubsub) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuitePubsub) TestNewPubsub() {
	assert.NotNil(this.T(), NewPubsub())
}

func (this *SuitePubsub) TestPubsub() {
	target := NewPubsub()
	valid := "param"
	count := 0
	target.Sub("proc1", func(param any) {
		if param.(string) == valid {
			count++
		} // if
	})
	target.Sub("proc1", func(param any) {
		if param.(string) == valid {
			count++
		} // if
	})
	target.Sub("proc2", func(param any) {
		if param.(string) == valid {
			count++
		} // if
	})
	target.Sub("proc2", func(param any) {
		if param.(string) == valid {
			count++
		} // if
	})
	target.Pub("proc1", valid)
	target.Pub("proc2", valid)
	assert.Equal(this.T(), 4, count)
}
