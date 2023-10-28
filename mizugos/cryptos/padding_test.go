package cryptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/testdata"
)

func TestPadding(t *testing.T) {
	suite.Run(t, new(SuitePadding))
}

type SuitePadding struct {
	suite.Suite
	testdata.Env
}

func (this *SuitePadding) SetupSuite() {
	this.Env = testdata.EnvSetup("test-cryptos-padding")
}

func (this *SuitePadding) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuitePadding) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuitePadding) TestPadding() {
	blockSize := 64
	length1 := 99
	length2 := 199
	source := []byte(helps.RandString(length1, helps.StrNumberAlpha))
	padstr := pad(PaddingZero, source, blockSize)
	result := unpad(PaddingZero, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingZero")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	source = []byte(helps.RandString(length2, helps.StrNumberAlpha))
	padstr = pad(PaddingZero, source, blockSize)
	result = unpad(PaddingZero, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingZero")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	source = []byte(helps.RandString(length1, helps.StrNumberAlpha))
	padstr = pad(PaddingPKCS7, source, blockSize)
	result = unpad(PaddingPKCS7, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingPKCS7")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	source = []byte(helps.RandString(length2, helps.StrNumberAlpha))
	padstr = pad(PaddingPKCS7, source, blockSize)
	result = unpad(PaddingPKCS7, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingPKCS7")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	assert.Nil(this.T(), pad(-1, nil, 0))
	assert.Nil(this.T(), unpad(-1, nil))
}
