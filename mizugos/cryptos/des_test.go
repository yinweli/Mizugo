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

func TestDes(t *testing.T) {
	suite.Run(t, new(SuiteDes))
}

type SuiteDes struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteDes) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-des"))
}

func (this *SuiteDes) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteDes) TestDesECB() {
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

func (this *SuiteDes) TestDesCBC() {
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

func (this *SuiteDes) TestRandDesKey() {
	key1 := RandDesKey()
	assert.NotNil(this.T(), key1)
	assert.Len(this.T(), key1, DesKeySize)

	key2 := RandDesKeyString()
	assert.NotNil(this.T(), key2)
	assert.Len(this.T(), key2, DesKeySize)
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
