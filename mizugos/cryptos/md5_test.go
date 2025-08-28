package cryptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestMD5(t *testing.T) {
	suite.Run(t, new(SuiteMD5))
}

type SuiteMD5 struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMD5) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-md5"))
}

func (this *SuiteMD5) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMD5) TestMD5String() {
	input := "test md5 string"
	output := MD5String([]byte(input))
	this.NotNil(output)
	fmt.Printf("%v => %v\n", input, output)
}

func BenchmarkMD5String1024(b *testing.B) {
	input := []byte(helps.RandString(1024, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkMD5String2048(b *testing.B) {
	input := []byte(helps.RandString(2048, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}

func BenchmarkMD5String4096(b *testing.B) {
	input := []byte(helps.RandString(4096, helps.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_ = MD5String(input)
	} // for
}
