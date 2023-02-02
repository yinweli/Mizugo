package cryptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestPadding(t *testing.T) {
	suite.Run(t, new(SuitePadding))
}

type SuitePadding struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	blockSize int
	length1   int
	length2   int
}

func (this *SuitePadding) SetupSuite() {
	this.Change("test-cryptos-padding")
	this.blockSize = 64
	this.length1 = 99
	this.length2 = 199
}

func (this *SuitePadding) TearDownSuite() {
	this.Restore()
}

func (this *SuitePadding) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuitePadding) TestPadding() {
	source := []byte(utils.RandString(this.length1))
	padstr := pad(PaddingZero, source, this.blockSize)
	result := unpad(PaddingZero, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingZero")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	source = []byte(utils.RandString(this.length2))
	padstr = pad(PaddingZero, source, this.blockSize)
	result = unpad(PaddingZero, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingZero")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	source = []byte(utils.RandString(this.length1))
	padstr = pad(PaddingPKCS7, source, this.blockSize)
	result = unpad(PaddingPKCS7, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingPKCS7")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	source = []byte(utils.RandString(this.length2))
	padstr = pad(PaddingPKCS7, source, this.blockSize)
	result = unpad(PaddingPKCS7, padstr)
	assert.Equal(this.T(), source, result)
	fmt.Println(">> PaddingPKCS7")
	fmt.Printf("source=%v\n", source)
	fmt.Printf("padstr=%v\n", padstr)
	fmt.Printf("result=%v\n", result)

	assert.Nil(this.T(), pad(-1, nil, 0))
	assert.Nil(this.T(), unpad(-1, nil))
}
