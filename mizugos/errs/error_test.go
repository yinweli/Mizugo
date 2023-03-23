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
	testdata.Env
}

func (this *SuiteErr) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-errs-error")
}

func (this *SuiteErr) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteErr) TearDownTest() {
	testdata.Leak(this.T(), true)
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
