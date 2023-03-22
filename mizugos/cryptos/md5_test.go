package cryptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMD5(t *testing.T) {
	suite.Run(t, new(SuiteMD5))
}

type SuiteMD5 struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteMD5) SetupSuite() {
	this.TBegin("test-cryptos-md5", "")
}

func (this *SuiteMD5) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteMD5) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteMD5) TestMD5String() {
	input := "test md5 string"
	output := MD5String([]byte(input))
	assert.NotNil(this.T(), output)
	fmt.Printf("%v => %v\n", input, output)
}

func BenchmarkMD5String1024(b *testing.B) {
	input := []byte(utils.RandString(1024))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkMD5String2048(b *testing.B) {
	input := []byte(utils.RandString(2048))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkMD5String4096(b *testing.B) {
	input := []byte(utils.RandString(4096))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}
