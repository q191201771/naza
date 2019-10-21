package nazamd5

import (
	"crypto/md5"
	"encoding/hex"
)

// 返回32字节小写字符串
func MD5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
