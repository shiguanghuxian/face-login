package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

/* 加密算法 */

// GetMd5String 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetMd5 生产md5字节
func GetMd5(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

// GetSha1String 对字符串进行SHA1哈希
func GetSha1String(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// UserPwdEncrypt 加密密码明文
func UserPwdEncrypt(password string, salts ...string) string {
	salt := "vivi_"
	if len(salts) > 0 {
		salt = salts[0]
	}
	return GetSha1String(GetSha1String(string(GetMd5(password))+salt) + GetMd5String(salt))
}

// HmacSha1ToString HmacSha1
func HmacSha1ToString(k, v string) string {
	key := []byte(v)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(k))
	return hex.EncodeToString(mac.Sum(nil))
}
