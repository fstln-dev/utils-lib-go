package encrypt_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fstln-dev/utils-lib-go/pkg/encrypt"
)

// TestHS256IsStandardHMAC 锁定 HS256 的算法，以及「secret 不得出现在输出里」。
//
// 旧实现写成 hex.EncodeToString(s.Sum([]byte(secret)))，而 hash.Sum(b) 的语义是
// 「把当前摘要追加到 b 之后」，输出因此是 hex(secret ‖ sha256(data))——前
// 2*len(secret) 个 hex 字符即 secret 明文。
func TestHS256IsStandardHMAC(t *testing.T) {

	secret := "zd89gu"
	parts := []string{"header", ".", "payload"}

	got, err := encrypt.HS256(secret, parts...)
	if err != nil {
		t.Fatalf("HS256: %s", err)
	}

	// 1) sha256 的 hex 长度固定 64；旧实现是 2*(len(secret)+32)。
	if len(got) != 64 {
		t.Errorf("输出长度 = %d，期望 64（sha256 hex）", len(got))
	}

	// 2) secret 不得以 hex 形式出现在输出里——旧实现会命中这一条。
	if strings.Contains(got, hex.EncodeToString([]byte(secret))) {
		t.Errorf("输出中泄漏了 secret：%s", got)
	}

	// 3) 与独立计算的 HMAC-SHA256 比对。
	mac := hmac.New(sha256.New, []byte(secret))
	for _, p := range parts {
		mac.Write([]byte(p))
	}
	if want := hex.EncodeToString(mac.Sum(nil)); got != want {
		t.Errorf("HS256 = %s，期望标准 HMAC-SHA256 = %s", got, want)
	}

	// 4) 换 secret 必须换输出。
	other, err := encrypt.HS256("other-secret", parts...)
	if err != nil {
		t.Fatalf("HS256: %s", err)
	}
	if other == got {
		t.Error("不同 secret 产生了相同输出")
	}
}
