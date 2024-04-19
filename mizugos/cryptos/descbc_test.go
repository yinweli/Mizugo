package cryptos

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestDesCBC(t *testing.T) {
	suite.Run(t, new(SuiteDesCBC))
}

type SuiteDesCBC struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteDesCBC) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-descbc"))
}

func (this *SuiteDesCBC) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteDesCBC) TestDesCBC() {
	keyValid := RandDesKeyString()
	keyInvalid := helps.RandString(15, helps.StrNumberAlpha)
	ivValid := RandDesKeyString()
	ivInvalid := ""
	target := NewDesCBC(PaddingPKCS7, keyValid, ivValid)
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

	_, err := NewDesCBC(PaddingPKCS7, keyInvalid, ivValid).Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivInvalid).Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivValid).Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivValid).Encode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivValid).Encode([]byte{})
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyInvalid, ivValid).Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivInvalid).Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivValid).Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivValid).Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = NewDesCBC(PaddingPKCS7, keyValid, ivValid).Decode([]byte{})
	assert.NotNil(this.T(), err)
}

func BenchmarkDesCBCEncode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))
	target := NewDesCBC(PaddingPKCS7, RandDesKeyString(), RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkDesCBCEncode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))
	target := NewDesCBC(PaddingPKCS7, RandDesKeyString(), RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkDesCBCEncode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))
	target := NewDesCBC(PaddingPKCS7, RandDesKeyString(), RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkDesCBCDecode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))
	target := NewDesCBC(PaddingPKCS7, RandDesKeyString(), RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}

func BenchmarkDesCBCDecode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))
	target := NewDesCBC(PaddingPKCS7, RandDesKeyString(), RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}

func BenchmarkDesCBCDecode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))
	target := NewDesCBC(PaddingPKCS7, RandDesKeyString(), RandDesKeyString())

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}
