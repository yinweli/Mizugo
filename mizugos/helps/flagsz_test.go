package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestFlagsz(t *testing.T) {
	suite.Run(t, new(SuiteFlagsz))
}

type SuiteFlagsz struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteFlagsz) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-flagsz"))
}

func (this *SuiteFlagsz) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteFlagsz) TestFlagsz() {
	target := FlagszInit(9, true)
	assert.Len(this.T(), target, 9)
	assert.True(this.T(), FlagszAny(target))
	assert.True(this.T(), FlagszAll(target))
	assert.False(this.T(), FlagszNone(target))
	assert.Equal(this.T(), int32(9), FlagszCount(target, true))
	assert.Equal(this.T(), int32(0), FlagszCount(target, false))

	target = FlagszInit(9, false)
	assert.Len(this.T(), target, 9)
	assert.False(this.T(), FlagszAny(target))
	assert.False(this.T(), FlagszAll(target))
	assert.True(this.T(), FlagszNone(target))
	assert.Equal(this.T(), int32(0), FlagszCount(target, true))
	assert.Equal(this.T(), int32(9), FlagszCount(target, false))

	target = FlagszAdd(target, true)
	assert.Len(this.T(), target, 10)
	assert.True(this.T(), FlagszGet(target, 9))

	target = FlagszSet(target, 10, false)
	assert.Len(this.T(), target, 11)
	assert.False(this.T(), FlagszGet(target, 10))

	assert.Equal(this.T(), "1000", FlagszAND("1100", "1010"))
	assert.Equal(this.T(), "1000", FlagszAND("110", "1010"))
	assert.Equal(this.T(), "1110", FlagszOR("1100", "1010"))
	assert.Equal(this.T(), "1110", FlagszOR("110", "1010"))
	assert.Equal(this.T(), "0110", FlagszXOR("1100", "1010"))
	assert.Equal(this.T(), "0110", FlagszXOR("110", "1010"))
}
