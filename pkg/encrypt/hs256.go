package encrypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HS256 以 secret 为密钥，对 data 顺序拼接后计算 HMAC-SHA256 并返回 hex。
//
// 旧实现写成 s.Sum([]byte(secret))，但 hash.Sum(b) 的语义是「把当前摘要追加到 b
// 之后」，产出的是 hex(secret ‖ sha256(data))——前 2*len(secret) 个 hex 字符即
// secret 明文。
func HS256(secret string, data ...string) (string, error) {

	mac := hmac.New(sha256.New, []byte(secret))

	for _, d := range data {
		if _, err := mac.Write([]byte(d)); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}
