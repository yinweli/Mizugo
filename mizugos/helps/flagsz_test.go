package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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

func (this *SuiteFlagsz) TestFlagszInit() {
	target := FlagszInit(2, true)
	this.Len(target, 2)
	this.True(FlagszAll(target))

	target = FlagszInit(2, false)
	this.Len(target, 2)
	this.True(FlagszNone(target))

	target = FlagszInit(-1, false)
	this.Empty(target)
}

func (this *SuiteFlagsz) TestFlagszSet() {
	target := FlagszSet("", 1, true)
	this.Len(target, 2)
	this.False(FlagszGet(target, 0))
	this.True(FlagszGet(target, 1))

	target = FlagszSet("", -1, true)
	this.Empty(target)
}

func (this *SuiteFlagsz) TestFlagszAdd() {
	target := FlagszAdd("", true)
	this.Len(target, 1)
	this.True(FlagszGet(target, 0))
}

func (this *SuiteFlagsz) TestFlagszAND() {
	target := FlagszAND("010", "101")
	this.Len(target, 3)
	this.True(FlagszNone(target))
}

func (this *SuiteFlagsz) TestFlagszOR() {
	target := FlagszOR("010", "101")
	this.Len(target, 3)
	this.True(FlagszAll(target))
}

func (this *SuiteFlagsz) TestFlagszXOR() {
	target := FlagszXOR("010", "111")
	this.Len(target, 3)
	this.True(FlagszGet(target, 0))
	this.False(FlagszGet(target, 1))
	this.True(FlagszGet(target, 2))
}

func (this *SuiteFlagsz) TestFlagszGet() {
	target := FlagszSet("", 1, true)
	this.False(FlagszGet(target, 0))
	this.True(FlagszGet(target, 1))
	this.False(FlagszGet(target, 2))
}

func (this *SuiteFlagsz) TestFlagszAny() {
	target := FlagszSet("", 1, true)
	this.True(FlagszAny(target))
	target = FlagszSet("", 1, false)
	this.False(FlagszAny(target))
}

func (this *SuiteFlagsz) TestFlagszAll() {
	target := FlagszInit(1, true)
	this.True(FlagszAll(target))
	target = FlagszInit(1, false)
	this.False(FlagszAll(target))
}

func (this *SuiteFlagsz) TestFlagszNone() {
	target := FlagszInit(1, true)
	this.False(FlagszNone(target))
	target = FlagszInit(1, false)
	this.True(FlagszNone(target))
}

func (this *SuiteFlagsz) TestFlagszCount() {
	target := FlagszSet("", 1, true)
	this.Equal(int32(1), FlagszCount(target, true))
	this.Equal(int32(1), FlagszCount(target, false))
}
