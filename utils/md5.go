package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(text string) string {
	harsher := md5.New()
	harsher.Write([]byte(text))
	return hex.EncodeToString(harsher.Sum(nil))
}
