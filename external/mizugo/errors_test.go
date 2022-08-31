package mizugo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestErrors(t *testing.T) {
	suite.Run(t, new(SuiteErrors))
}

type SuiteErrors struct {
	suite.Suite
	errID      ErrID
	errAlone   error
	errInside  error
	errOutside error
}

func (this *SuiteErrors) SetupSuite() {
	this.errID = 1234
	this.errAlone = fmt.Errorf("alone")
	this.errInside = fmt.Errorf("inside")
	this.errOutside = fmt.Errorf("outside: %w", this.errInside)
}

func (this *SuiteErrors) target() error {
	return Errorf(this.errID, this.errOutside)
}

func (this *SuiteErrors) TestErrorf() {
	target := this.target()

	assert.NotNil(this.T(), target)
	assert.True(this.T(), errors.Is(target, this.errInside))
	assert.True(this.T(), errors.Is(target, this.errOutside))
	assert.False(this.T(), errors.Is(target, this.errAlone))
}

func (this *SuiteErrors) TestUnwrapErrID() {
	target := this.target()

	assert.Equal(this.T(), this.errID, UnwrapErrID(target))
	assert.Equal(this.T(), UnknownErrID, UnwrapErrID(this.errInside))
	assert.Equal(this.T(), UnknownErrID, UnwrapErrID(this.errOutside))
}

func TestWrapError(t *testing.T) {
	suite.Run(t, new(SuiteWrapError))
}

type SuiteWrapError struct {
	suite.Suite
	errID ErrID
	err   error
}

func (this *SuiteWrapError) SetupSuite() {
	this.errID = 1234
	this.err = fmt.Errorf("err")
}

func (this *SuiteWrapError) target() wrapError {
	return wrapError{errID: this.errID, err: this.err}
}

func (this *SuiteWrapError) TestErrID() {
	target := this.target()

	assert.Equal(this.T(), this.errID, target.ErrID())
}

func (this *SuiteWrapError) TestUnwrap() {
	target := this.target()

	assert.Equal(this.T(), this.err, target.Unwrap())
}

func (this *SuiteWrapError) TestError() {
	target := this.target()

	assert.NotNil(this.T(), target.Error())
}
