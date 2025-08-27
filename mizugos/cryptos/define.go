package cryptos

import (
	"github.com/yinweli/Mizugo/v2/mizugos/helps"
)

const DesKeySize = 8                      // des密鑰長度
const DesKeyLetter = helps.StrNumberAlpha // des密鑰字串

// RandDesKey 產生隨機des密鑰
func RandDesKey() []byte {
	return []byte(helps.RandString(DesKeySize, DesKeyLetter))
}

// RandDesKeyString 產生隨機des密鑰字串
func RandDesKeyString() string {
	return string(RandDesKey())
}
