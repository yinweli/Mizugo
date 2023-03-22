package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestErr(t *testing.T) {
	suite.Run(t, new(SuiteErr))
}

type SuiteErr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteErr) SetupSuite() {
	this.TBegin("test-errs-error", "")
}

func (this *SuiteErr) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteErr) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteErr) TestError() {
	err := Errort("tag")
	assert.NotNil(this.T(), err)
	fmt.Println(err)

	err = Errore("tag", fmt.Errorf("errore: %v", 1))
	assert.NotNil(this.T(), err)
	fmt.Println(err)

	err = Errorf("tag", "errorf: %v", 1)
	assert.NotNil(this.T(), err)
	fmt.Println(err)
}

func (this *SuiteErr) TestWrapError() {
	err := &wrapError{
		tag: "tag",
		err: fmt.Errorf("error"),
	}
	assert.NotNil(this.T(), err.Error())
	fmt.Println(err.Error())

	err.tag = nil
	assert.NotNil(this.T(), err.Error())
	fmt.Println(err.Error())

	err.err = nil
	assert.NotNil(this.T(), err.Error())
	fmt.Println(err.Error())
}
