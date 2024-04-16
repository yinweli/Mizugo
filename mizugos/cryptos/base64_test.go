package cryptos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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
	packet := Base64Encode(input)
	output, err := Base64Decode(packet)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), input, output)
}

func BenchmarkBase64Encode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Encode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Encode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Decode1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}

func BenchmarkBase64Decode2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}

func BenchmarkBase64Decode4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}
