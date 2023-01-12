package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEncoding(t *testing.T) {
	suite.Run(t, new(SuiteEncoding))
}

type SuiteEncoding struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteEncoding) SetupSuite() {
	this.Change("test-utils-encoding")
}

func (this *SuiteEncoding) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEncoding) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteEncoding) TestMD5String() {
	input := "test md5 string"
	output := MD5String([]byte(input))
	assert.NotNil(this.T(), output)
	fmt.Printf("%v => %v\n", input, output)
}

func (this *SuiteEncoding) TestBase64() {
	input := []byte("test base64 string")
	packet := Base64Encode(input)
	output, err := Base64Decode(packet)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), input, output)
}

func BenchmarkMD5String1024(b *testing.B) {
	input := []byte(RandString(1024))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkMD5String2048(b *testing.B) {
	input := []byte(RandString(2048))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkMD5String4096(b *testing.B) {
	input := []byte(RandString(4096))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkBase64Encode1024(b *testing.B) {
	input := []byte(RandString(1024))

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Encode2048(b *testing.B) {
	input := []byte(RandString(2048))

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Encode4096(b *testing.B) {
	input := []byte(RandString(4096))

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Decode1024(b *testing.B) {
	input := []byte(RandString(1024))

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}

func BenchmarkBase64Decode2048(b *testing.B) {
	input := []byte(RandString(2048))

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}

func BenchmarkBase64Decode4096(b *testing.B) {
	input := []byte(RandString(4096))

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}
