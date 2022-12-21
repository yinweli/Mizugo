package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEncoding(t *testing.T) {
	suite.Run(t, new(SuiteEncoding))
}

type SuiteEncoding struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEncoding) SetupSuite() {
	this.Change("test-utils-encoding")
}

func (this *SuiteEncoding) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEncoding) TearDownTest() {
	goleak.VerifyNone(this.T())
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

func BenchmarkMD5String(b *testing.B) {
	input := []byte("benchmark md5 string")

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkBase64Encode(b *testing.B) {
	input := []byte("benchmark base64 encode string")

	for i := 0; i < b.N; i++ {
		_ = Base64Encode(input)
	} // for
}

func BenchmarkBase64Decode(b *testing.B) {
	input := []byte("benchmark base64 decode string")

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(input)
	} // for
}
