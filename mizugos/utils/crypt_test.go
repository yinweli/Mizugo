package utils

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestCrypt(t *testing.T) {
	suite.Run(t, new(SuiteCrypt))
}

type SuiteCrypt struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteCrypt) SetupSuite() {
	this.Change("test-utils-crypt")
}

func (this *SuiteCrypt) TearDownSuite() {
	this.Restore()
}

func (this *SuiteCrypt) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteCrypt) TestDesEncryptDecrypt() {
	for size := 16; size <= 64; size++ {
		key := RandDesKey()
		input := []byte(RandString(size))
		crypt, err := DesEncrypt(key, input)
		assert.Nil(this.T(), err)
		output, err := DesDecrypt(key, crypt)
		assert.Nil(this.T(), err)
		assert.Equal(this.T(), input, output)

		fmt.Printf("----- size: %v -----\n", size)
		fmt.Printf("key=%v\n", string(key))
		fmt.Printf("input=%v\n", string(input))
		fmt.Printf("output=%v\n", string(output))
		fmt.Printf("crypt=%v\n", hex.EncodeToString(crypt))
	} // for

	key := []byte(RandString(15))
	data := []byte(RandString(15))

	_, err := DesEncrypt(key, data)
	assert.NotNil(this.T(), err)

	_, err = DesDecrypt(key, data)
	assert.NotNil(this.T(), err)
}

func (this *SuiteCrypt) TestRandDesKey() {
	key1 := RandDesKey()
	assert.NotNil(this.T(), key1)
	assert.Len(this.T(), key1, DesKeySize)

	key2 := RandDesKeyString()
	assert.NotNil(this.T(), key2)
	assert.Len(this.T(), key2, DesKeySize)
}

func BenchmarkDesEncrypt1024(b *testing.B) {
	key := RandDesKey()
	input := []byte(RandString(1024))

	for i := 0; i < b.N; i++ {
		_, _ = DesEncrypt(key, input)
	} // for
}

func BenchmarkDesEncrypt2048(b *testing.B) {
	key := RandDesKey()
	input := []byte(RandString(2048))

	for i := 0; i < b.N; i++ {
		_, _ = DesEncrypt(key, input)
	} // for
}

func BenchmarkDesEncrypt4096(b *testing.B) {
	key := RandDesKey()
	input := []byte(RandString(4096))

	for i := 0; i < b.N; i++ {
		_, _ = DesEncrypt(key, input)
	} // for
}

func BenchmarkDesDecrypt1024(b *testing.B) {
	key := RandDesKey()
	input := []byte(RandString(1024))

	for i := 0; i < b.N; i++ {
		_, _ = DesDecrypt(key, input)
	} // for
}

func BenchmarkDesDecrypt2048(b *testing.B) {
	key := RandDesKey()
	input := []byte(RandString(2048))

	for i := 0; i < b.N; i++ {
		_, _ = DesDecrypt(key, input)
	} // for
}

func BenchmarkDesDecrypt4096(b *testing.B) {
	key := RandDesKey()
	input := []byte(RandString(4096))

	for i := 0; i < b.N; i++ {
		_, _ = DesDecrypt(key, input)
	} // for
}
