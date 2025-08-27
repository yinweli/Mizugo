package cryptos

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestDesECB(t *testing.T) {
	suite.Run(t, new(SuiteDesECB))
}

type SuiteDesECB struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteDesECB) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-desecb"))
}

func (this *SuiteDesECB) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteDesECB) TestDesECB() {
	keyValid := RandDesKeyString()
	keyInvalid := helps.RandString(15, helps.StrNumberAlpha)
	target := NewDesECB(PaddingPKCS7, keyValid)
	assert.NotNil(this.T(), target)

	for size := 16; size <= 64; size++ {
		input := []byte(helps.RandString(size, helps.StrNumberAlpha))
		crypto, err := target.Encode(input)
		assert.Nil(this.T(), err)
		output, err := target.Decode(crypto)
		assert.Nil(this.T(), err)
		assert.Equal(this.T(), input, output)

		fmt.Printf("----- size: %v -----\n", size)
		fmt.Printf("input=%v\n", string(input))
		fmt.Printf("output=%v\n", string(output.([]byte)))
		fmt.Printf("crypto=%v\n", hex.EncodeToString(crypto.([]byte)))
	} // for

	_, err := NewDesECB(PaddingPKCS7, keyInvalid).Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyValid).Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyValid).Encode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyValid).Encode([]byte{})
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyInvalid).Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyValid).Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyValid).Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = NewDesECB(PaddingPKCS7, keyValid).Decode([]byte{})
	assert.NotNil(this.T(), err)
}

func BenchmarkDesECBEncode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))
	target := NewDesECB(PaddingPKCS7, RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkDesECBEncode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))
	target := NewDesECB(PaddingPKCS7, RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkDesECBEncode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))
	target := NewDesECB(PaddingPKCS7, RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkDesECBDecode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))
	target := NewDesECB(PaddingPKCS7, RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}

func BenchmarkDesECBDecode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))
	target := NewDesECB(PaddingPKCS7, RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}

func BenchmarkDesECBDecode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))
	target := NewDesECB(PaddingPKCS7, RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}
