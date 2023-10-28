package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestFit(t *testing.T) {
	suite.Run(t, new(SuiteFit))
}

type SuiteFit struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteFit) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-fit")
}

func (this *SuiteFit) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteFit) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteFit) TestFit() {
	target := NewFit[int](
		func() int { return 10 },
		func() int { return 5 },
	)
	assert.NotNil(this.T(), target)
	result, remain, added := target.Check(7)
	assert.Equal(this.T(), 7, result)
	assert.Equal(this.T(), 0, remain)
	assert.Equal(this.T(), 0, added)
	result, remain, added = target.Check(11)
	assert.Equal(this.T(), 10, result)
	assert.Equal(this.T(), 1, remain)
	assert.Equal(this.T(), -1, added)
	result, remain, added = target.Check(3)
	assert.Equal(this.T(), 5, result)
	assert.Equal(this.T(), -2, remain)
	assert.Equal(this.T(), 2, added)
	result, remain, added = target.Check(5, 1)
	assert.Equal(this.T(), 6, result)
	assert.Equal(this.T(), 0, remain)
	assert.Equal(this.T(), 1, added)
	result, remain, added = target.Check(10, -1)
	assert.Equal(this.T(), 9, result)
	assert.Equal(this.T(), 0, remain)
	assert.Equal(this.T(), -1, added)
	result, remain, added = target.Check(5, 10)
	assert.Equal(this.T(), 10, result)
	assert.Equal(this.T(), 5, remain)
	assert.Equal(this.T(), 5, added)
	result, remain, added = target.Check(10, -10)
	assert.Equal(this.T(), 5, result)
	assert.Equal(this.T(), -5, remain)
	assert.Equal(this.T(), -5, added)

	target = NewFit[int](
		func() int { return 10 },
		nil,
	)
	result, remain, added = target.Check(7)
	assert.Equal(this.T(), 7, result)
	assert.Equal(this.T(), 0, remain)
	assert.Equal(this.T(), 0, added)
	result, remain, added = target.Check(11)
	assert.Equal(this.T(), 10, result)
	assert.Equal(this.T(), 1, remain)
	assert.Equal(this.T(), -1, added)
	result, remain, added = target.Check(-2)
	assert.Equal(this.T(), 0, result)
	assert.Equal(this.T(), -2, remain)
	assert.Equal(this.T(), 2, added)

	target = NewFit[int](
		nil,
		func() int { return -10 },
	)
	result, remain, added = target.Check(-5)
	assert.Equal(this.T(), -5, result)
	assert.Equal(this.T(), 0, remain)
	assert.Equal(this.T(), 0, added)
	result, remain, added = target.Check(5)
	assert.Equal(this.T(), 0, result)
	assert.Equal(this.T(), 5, remain)
	assert.Equal(this.T(), -5, added)
	result, remain, added = target.Check(-12)
	assert.Equal(this.T(), -10, result)
	assert.Equal(this.T(), -2, remain)
	assert.Equal(this.T(), 2, added)

	target = NewFit[int](nil, nil)
	result, remain, added = target.Check(0)
	assert.Equal(this.T(), 0, result)
	assert.Equal(this.T(), 0, remain)
	assert.Equal(this.T(), 0, added)
	result, remain, added = target.Check(2)
	assert.Equal(this.T(), 0, result)
	assert.Equal(this.T(), 2, remain)
	assert.Equal(this.T(), -2, added)
	result, remain, added = target.Check(-2)
	assert.Equal(this.T(), 0, result)
	assert.Equal(this.T(), -2, remain)
	assert.Equal(this.T(), 2, added)
}
