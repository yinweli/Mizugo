package cryptos

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestDes(t *testing.T) {
	suite.Run(t, new(SuiteDes))
}

type SuiteDes struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteDes) SetupSuite() {
	this.Env = testdata.EnvSetup("test-cryptos-des")
}

func (this *SuiteDes) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteDes) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteDes) TestDesECB() {
	for size := 16; size <= 64; size++ {
		key := RandDesKey()
		input := []byte(utils.RandString(size, utils.StrNumberAlpha))
		crypto, err := DesECBEncrypt(PaddingPKCS7, key, input)
		assert.Nil(this.T(), err)
		output, err := DesECBDecrypt(PaddingPKCS7, key, crypto)
		assert.Nil(this.T(), err)
		assert.Equal(this.T(), input, output)

		fmt.Printf("----- size: %v -----\n", size)
		fmt.Printf("key=%v\n", string(key))
		fmt.Printf("input=%v\n", string(input))
		fmt.Printf("output=%v\n", string(output))
		fmt.Printf("crypto=%v\n", hex.EncodeToString(crypto))
	} // for

	keyValid := []byte(utils.RandString(DesKeySize, utils.StrNumberAlpha))
	keyInvalid := []byte(utils.RandString(15, utils.StrNumberAlpha))
	dataInvalid := []byte(utils.RandString(15, utils.StrNumberAlpha))

	_, err := DesECBEncrypt(PaddingPKCS7, keyInvalid, dataInvalid)
	assert.NotNil(this.T(), err)

	_, err = DesECBDecrypt(PaddingPKCS7, keyInvalid, dataInvalid)
	assert.NotNil(this.T(), err)

	_, err = DesECBDecrypt(PaddingPKCS7, keyValid, dataInvalid)
	assert.NotNil(this.T(), err)
}

func (this *SuiteDes) TestDesCBC() {
	for size := 16; size <= 64; size++ {
		key := RandDesKey()
		input := []byte(utils.RandString(size, utils.StrNumberAlpha))
		crypto, err := DesCBCEncrypt(PaddingPKCS7, key, key, input)
		assert.Nil(this.T(), err)
		output, err := DesCBCDecrypt(PaddingPKCS7, key, key, crypto)
		assert.Nil(this.T(), err)
		assert.Equal(this.T(), input, output)

		fmt.Printf("----- size: %v -----\n", size)
		fmt.Printf("key=%v\n", string(key))
		fmt.Printf("input=%v\n", string(input))
		fmt.Printf("output=%v\n", string(output))
		fmt.Printf("crypto=%v\n", hex.EncodeToString(crypto))
	} // for

	keyValid := []byte(utils.RandString(DesKeySize, utils.StrNumberAlpha))
	keyInvalid := []byte(utils.RandString(15, utils.StrNumberAlpha))
	dataInvalid := []byte(utils.RandString(15, utils.StrNumberAlpha))

	_, err := DesCBCEncrypt(PaddingPKCS7, keyInvalid, keyValid, dataInvalid)
	assert.NotNil(this.T(), err)

	_, err = DesCBCEncrypt(PaddingPKCS7, keyValid, keyInvalid, dataInvalid)
	assert.NotNil(this.T(), err)

	_, err = DesCBCDecrypt(PaddingPKCS7, keyInvalid, keyValid, dataInvalid)
	assert.NotNil(this.T(), err)

	_, err = DesCBCDecrypt(PaddingPKCS7, keyValid, keyInvalid, dataInvalid)
	assert.NotNil(this.T(), err)

	_, err = DesCBCDecrypt(PaddingPKCS7, keyValid, keyValid, dataInvalid)
	assert.NotNil(this.T(), err)
}

func (this *SuiteDes) TestRandDesKey() {
	key1 := RandDesKey()
	assert.NotNil(this.T(), key1)
	assert.Len(this.T(), key1, DesKeySize)

	key2 := RandDesKeyString()
	assert.NotNil(this.T(), key2)
	assert.Len(this.T(), key2, DesKeySize)
}

func BenchmarkDesECBEncrypt1024(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(1024, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesECBEncrypt(PaddingPKCS7, key, input)
	} // for
}

func BenchmarkDesECBEncrypt2048(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(2048, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesECBEncrypt(PaddingPKCS7, key, input)
	} // for
}

func BenchmarkDesECBEncrypt4096(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(4096, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesECBEncrypt(PaddingPKCS7, key, input)
	} // for
}

func BenchmarkDesECBDecrypt1024(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(1024, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesECBDecrypt(PaddingPKCS7, key, input)
	} // for
}

func BenchmarkDesECBDecrypt2048(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(2048, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesECBDecrypt(PaddingPKCS7, key, input)
	} // for
}

func BenchmarkDesECBDecrypt4096(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(4096, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesECBDecrypt(PaddingPKCS7, key, input)
	} // for
}

func BenchmarkDesCBCEncrypt1024(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(1024, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesCBCEncrypt(PaddingPKCS7, key, key, input)
	} // for
}

func BenchmarkDesCBCEncrypt2048(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(2048, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesCBCEncrypt(PaddingPKCS7, key, key, input)
	} // for
}

func BenchmarkDesCBCEncrypt4096(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(4096, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesCBCEncrypt(PaddingPKCS7, key, key, input)
	} // for
}

func BenchmarkDesCBCDecrypt1024(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(1024, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesCBCDecrypt(PaddingPKCS7, key, key, input)
	} // for
}

func BenchmarkDesCBCDecrypt2048(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(2048, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesCBCDecrypt(PaddingPKCS7, key, key, input)
	} // for
}

func BenchmarkDesCBCDecrypt4096(b *testing.B) {
	key := RandDesKey()
	input := []byte(utils.RandString(4096, utils.StrNumberAlpha))

	for i := 0; i < b.N; i++ {
		_, _ = DesCBCDecrypt(PaddingPKCS7, key, key, input)
	} // for
}
