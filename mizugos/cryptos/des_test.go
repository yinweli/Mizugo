package cryptos

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestPadding(t *testing.T) {
	suite.Run(t, new(SuitePadding))
}

type SuitePadding struct {
	suite.Suite
	trials.Catalog
}

func (this *SuitePadding) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-padding"))
}

func (this *SuitePadding) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuitePadding) TestRandDesKey() {
	target := RandDesKey()
	this.NotNil(target)
	this.Len(target, DesKeySize)
}

func (this *SuitePadding) TestRandDesKeyString() {
	target := RandDesKeyString()
	this.NotNil(target)
	this.Len(target, DesKeySize)
}

func (this *SuitePadding) TestPadding() {
	source := []byte(helps.RandString(99, helps.StrNumberAlpha))
	result, err := pad(PaddingZero, source, 64)
	this.Nil(err)
	result, err = unpad(PaddingZero, result)
	this.Nil(err)
	this.Equal(source, result)

	source = []byte(helps.RandString(99, helps.StrNumberAlpha))
	result, err = pad(PaddingPKCS7, source, 64)
	this.Nil(err)
	result, err = unpad(PaddingPKCS7, result)
	this.Nil(err)
	this.Equal(source, result)

	_, err = pad(-1, nil, 0)
	this.NotNil(err)
	_, err = unpad(-1, nil)
	this.NotNil(err)
}

func (this *SuitePadding) TestZeroPad() {
	source := []byte(helps.RandString(99, helps.StrNumberAlpha))
	result, err := zeroPad(source, 64)
	this.Nil(err)
	result, err = zeroUnpad(result)
	this.Nil(err)
	this.Equal(source, result)

	_, err = zeroPad(nil, 0)
	this.NotNil(err)
}

func (this *SuitePadding) TestPKCS7Pad() {
	source := []byte(helps.RandString(99, helps.StrNumberAlpha))
	result, err := pkcs7Pad(source, 64)
	this.Nil(err)
	result, err = pkcs7Unpad(result)
	this.Nil(err)
	this.Equal(source, result)

	source = []byte(helps.RandString(199, helps.StrNumberAlpha))
	result, err = pkcs7Pad(source, 64)
	this.Nil(err)
	result, err = pkcs7Unpad(result)
	this.Nil(err)
	this.Equal(source, result)

	_, err = pkcs7Pad(source, 0)
	this.NotNil(err)
	_, err = pkcs7Unpad(nil)
	this.NotNil(err)
	_, err = pkcs7Unpad([]byte{byte(0)})
	this.NotNil(err)
	_, err = pkcs7Unpad([]byte{byte(1), byte(2), byte(3)})
	this.NotNil(err)
}
