package packets

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5String 取得md5字串
func MD5String(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
