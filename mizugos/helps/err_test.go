package helps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestErr(t *testing.T) {
	suite.Run(t, new(SuiteErr))
}

type SuiteErr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteErr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-err"))
}

func (this *SuiteErr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteErr) TestErr() {
	err := Err(nil)
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err(Success, ErrUnknown)
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err("err", ErrUnknown)
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err("err", uint64(ErrUnknown))
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err(fmt.Errorf("err"), ErrUnknown)
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err(&Error{
		err:   "err",
		errID: Success,
	}, ErrUnknown)
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err("str1", "str2", fmt.Errorf("err"), ErrUnknown)
	assert.NotNil(this.T(), err)
	assert.NotEmpty(this.T(), err.Error())
	assert.Equal(this.T(), ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	assert.Equal(this.T(), ErrUnwrap, UnwrapErrID(nil))
}
