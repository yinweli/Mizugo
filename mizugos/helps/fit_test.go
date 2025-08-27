package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestFit(t *testing.T) {
	suite.Run(t, new(SuiteFit))
}

type SuiteFit struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteFit) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-fit"))
}

func (this *SuiteFit) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteFit) TestFit() {
	target := NewFit[int](
		func() int { return 10 },
		func() int { return 5 },
	)
	this.NotNil(target)
	result, remain, added := target.Check(7)
	this.Equal(7, result)
	this.Equal(0, remain)
	this.Equal(0, added)
	result, remain, added = target.Check(11)
	this.Equal(10, result)
	this.Equal(1, remain)
	this.Equal(-1, added)
	result, remain, added = target.Check(3)
	this.Equal(5, result)
	this.Equal(-2, remain)
	this.Equal(2, added)
	result, remain, added = target.Check(5, 1)
	this.Equal(6, result)
	this.Equal(0, remain)
	this.Equal(1, added)
	result, remain, added = target.Check(10, -1)
	this.Equal(9, result)
	this.Equal(0, remain)
	this.Equal(-1, added)
	result, remain, added = target.Check(5, 10)
	this.Equal(10, result)
	this.Equal(5, remain)
	this.Equal(5, added)
	result, remain, added = target.Check(10, -10)
	this.Equal(5, result)
	this.Equal(-5, remain)
	this.Equal(-5, added)

	target = NewFit[int](
		func() int { return 10 },
		nil,
	)
	result, remain, added = target.Check(7)
	this.Equal(7, result)
	this.Equal(0, remain)
	this.Equal(0, added)
	result, remain, added = target.Check(11)
	this.Equal(10, result)
	this.Equal(1, remain)
	this.Equal(-1, added)
	result, remain, added = target.Check(-2)
	this.Equal(0, result)
	this.Equal(-2, remain)
	this.Equal(2, added)

	target = NewFit[int](
		nil,
		func() int { return -10 },
	)
	result, remain, added = target.Check(-5)
	this.Equal(-5, result)
	this.Equal(0, remain)
	this.Equal(0, added)
	result, remain, added = target.Check(5)
	this.Equal(0, result)
	this.Equal(5, remain)
	this.Equal(-5, added)
	result, remain, added = target.Check(-12)
	this.Equal(-10, result)
	this.Equal(-2, remain)
	this.Equal(2, added)

	target = NewFit[int](nil, nil)
	result, remain, added = target.Check(0)
	this.Equal(0, result)
	this.Equal(0, remain)
	this.Equal(0, added)
	result, remain, added = target.Check(2)
	this.Equal(0, result)
	this.Equal(2, remain)
	this.Equal(-2, added)
	result, remain, added = target.Check(-2)
	this.Equal(0, result)
	this.Equal(-2, remain)
	this.Equal(2, added)
}
