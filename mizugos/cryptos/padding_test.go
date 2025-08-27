package cryptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestPadding(t *testing.T) {
	suite.Run(t, new(SuitePadding))
}

type SuitePadding struct {
	suite.Suite
	trials.Catalog
}

func (this *SuitePadding) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-padding"))
}

func (this *SuitePadding) TearDownSuite() {
	trials.Restore(this.Catalog)
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
