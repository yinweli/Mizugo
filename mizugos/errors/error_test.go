package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestError(t *testing.T) {
	suite.Run(t, new(SuiteError))
}

type SuiteError struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteError) SetupSuite() {
	this.Change("test-errors-error")
}

func (this *SuiteError) TearDownSuite() {
	this.Restore()
}

func (this *SuiteError) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteError) TestErrorf() {
	alone := fmt.Errorf("alone")
	inside := fmt.Errorf("inside")
	outside := fmt.Errorf("outside: %w", inside)
	errorID := ErrorID(1)
	target := Errorf(errorID, outside)

	assert.NotNil(this.T(), target)
	assert.True(this.T(), errors.Is(target, inside))
	assert.True(this.T(), errors.Is(target, outside))
	assert.False(this.T(), errors.Is(target, alone))
}

func (this *SuiteError) TestUnwrapErrorID() {
	alone := fmt.Errorf("alone")
	errorID := ErrorID(1)
	target := Errorf(errorID, alone)

	assert.Equal(this.T(), errorID, UnwrapErrorID(target))
	assert.Equal(this.T(), Unknown, UnwrapErrorID(alone))
}

func (this *SuiteError) TestWrapError() {
	target := wrapError{errorID: 1, err: fmt.Errorf("error")}

	assert.Equal(this.T(), target.errorID, target.ErrorID())
	assert.Equal(this.T(), target.err, target.Unwrap())
	assert.Equal(this.T(), fmt.Sprintf("[%v] %v", target.errorID, target.err.Error()), target.Error())
}
