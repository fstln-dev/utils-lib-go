package jwt_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/fstln-dev/utils-lib-go/pkg/jwt"
)

type payloadS struct {
	A string        `json:"d"`
	B int           `json:"b"`
	C time.Duration `json:"c"`
}

// TestHS256IsStandardHMAC 锁定签名算法的回归用例。
//
// 旧实现写成 hex.EncodeToString(s.Sum([]byte(secret)))，而 hash.Sum(b) 的语义是
// 「把当前摘要追加到 b 之后」，签名因此是 hex(secret ‖ sha256(data))——前
// 2*len(secret) 个 hex 字符即 secret 明文，任何拿到 token 的人都能读出 secret
// 并伪造该账号的 token。本用例同时锁定「算法为标准 HMAC-SHA256」与「secret
// 不得出现在签名里」两件事。
func TestHS256IsStandardHMAC(t *testing.T) {

	secret := "zd89gu"

	tok, err := jwt.JWTEncode(secret, &payloadS{A: "12dsaf4", B: 1122, C: 2 * time.Second})
	if err != nil {
		t.Fatalf("JWTEncode: %s", err)
	}

	ss := strings.Split(tok, ".")
	if len(ss) != 3 {
		t.Fatalf("token 段数 = %d，期望 3：%s", len(ss), tok)
	}
	sig := ss[2]

	// 1) sha256 的 hex 长度固定 64；旧实现是 2*(len(secret)+32)。
	if len(sig) != 64 {
		t.Errorf("签名长度 = %d，期望 64（sha256 hex）", len(sig))
	}

	// 2) secret 不得以 hex 形式出现在签名里——旧实现会命中这一条。
	if strings.Contains(sig, hex.EncodeToString([]byte(secret))) {
		t.Errorf("签名中泄漏了 secret：sig=%s", sig)
	}

	// 3) 与独立计算的 HMAC-SHA256 比对（不经过 jwtHS256）。
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ss[0] + "." + ss[1]))
	if want := hex.EncodeToString(mac.Sum(nil)); sig != want {
		t.Errorf("签名 = %s，期望标准 HMAC-SHA256 = %s", sig, want)
	}

	// 4) 正确 secret 通过，错误 secret 拒绝。
	if ok, err := jwt.JWTValid(secret, tok); err != nil || !ok {
		t.Errorf("正确 secret 验签失败：ok=%v err=%v", ok, err)
	}
	if ok, _ := jwt.JWTValid("wrong-secret", tok); ok {
		t.Error("错误 secret 竟然验签通过")
	}
}
