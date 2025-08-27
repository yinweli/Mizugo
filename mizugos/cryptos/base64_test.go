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

func TestBase64(t *testing.T) {
	suite.Run(t, new(SuiteBase64))
}

type SuiteBase64 struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteBase64) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-base64"))
}

func (this *SuiteBase64) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteBase64) TestBase64() {
	input := []byte("test base64 string")
	target := NewBase64()
	assert.NotNil(this.T(), target)
	crypto, err := target.Encode(input)
	assert.Nil(this.T(), err)
	output, err := target.Decode(crypto)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), input, output)
	fmt.Printf("input=%v\n", string(input))
	fmt.Printf("output=%v\n", string(output.([]byte)))
	fmt.Printf("crypto=%v\n", hex.EncodeToString(crypto.([]byte)))
}

func BenchmarkBase64Encode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))
	target := NewBase64()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkBase64Encode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))
	target := NewBase64()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkBase64Encode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))
	target := NewBase64()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkBase64Decode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))
	target := NewBase64()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}

func BenchmarkBase64Decode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))
	target := NewBase64()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}

func BenchmarkBase64Decode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))
	target := NewBase64()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(input)
	} // for
}
