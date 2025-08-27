package helps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err(Success, ErrUnknown)
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err("err", ErrUnknown)
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err("err", uint64(ErrUnknown))
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err(fmt.Errorf("err"), ErrUnknown)
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err(&Error{
		err:   "err",
		errID: Success,
	}, ErrUnknown)
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	err = Err("str1", "str2", fmt.Errorf("err"), ErrUnknown)
	this.NotNil(err)
	this.NotEmpty(err.Error())
	this.Equal(ErrUnknown, UnwrapErrID(err))
	fmt.Println(err)

	this.Equal(ErrUnwrap, UnwrapErrID(nil))
}
