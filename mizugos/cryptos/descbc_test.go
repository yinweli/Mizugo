package cryptos

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	key1 := RandDesKeyString()
	key2 := helps.RandString(15, helps.StrNumberAlpha)
	iv1 := RandDesKeyString()
	iv2 := ""
	target := NewDesCBC(PaddingPKCS7, key1, iv1)
	this.NotNil(target)

	for size := 16; size <= 64; size++ {
		input := []byte(helps.RandString(size, helps.StrNumberAlpha))
		crypto, err := target.Encode(input)
		this.Nil(err)
		output, err := target.Decode(crypto)
		this.Nil(err)
		this.Equal(input, output)

		fmt.Printf("size: %v\n", size)
		fmt.Printf("  input  %v\n", string(input))
		fmt.Printf("  output %v\n", string(output.([]byte)))
		fmt.Printf("  crypto %v\n", hex.EncodeToString(crypto.([]byte)))
	} // for

	_, err := NewDesCBC(PaddingPKCS7, key1, iv1).Encode(nil)
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key1, iv2).Encode(nil)
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key1, iv1).Encode(testdata.Unknown)
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key1, iv1).Encode([]byte{})
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key2, iv1).Encode(nil)
	this.NotNil(err)

	_, err = NewDesCBC(PaddingPKCS7, key1, iv1).Decode(nil)
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key1, iv2).Decode(nil)
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key1, iv1).Decode(testdata.Unknown)
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key1, iv1).Decode([]byte{})
	this.NotNil(err)
	_, err = NewDesCBC(PaddingPKCS7, key2, iv1).Decode(nil)
	this.NotNil(err)
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
