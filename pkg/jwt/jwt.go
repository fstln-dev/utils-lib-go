package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

/**
* 注意事项记录：
* 1. 需要服务器时间同步
* 2. 为了保证安全，注意防范以下情况
*   2.1 使用 https，防止中间人攻击
*   2.2 （IP 固定）防范 reply 重放攻击
*   2.3 前端防范 xss 攻击
*   2.4 可建立版本号黑名单
* 3. 权限服务需要建立刷新机制
**/

type Header func(h *header)
type header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

func JwtTyp(t string) Header {
	return Header(func(h *header) {
		h.Typ = t
	})
}
func JwtAlg(s string) Header {
	return Header(func(h *header) {
		h.Alg = s
	})
}

// JWTEncode is the method to sign and generate JWT
// JWT 生成 + 签名
func JWTEncode(secret string, m interface{}, heads ...Header) (string, error) {
	h := header{
		Typ: "JWT",
		Alg: "HS256",
	}

	for _, head := range heads {
		head(&h)
	}

	hB, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	data := fmt.Sprintf(
		"%v.%v",
		base64.URLEncoding.EncodeToString(hB),
		base64.URLEncoding.EncodeToString(bytes),
	)

	sign := ""
	if h.Alg == "HS256" {
		sign, err = jwtHS256(secret, data)
	} else if h.Alg == "RS256" {
		// sign = rsa
	}

	if sign == "" {
		return "", errors.New("sign no find")
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", data, sign), nil
}

// JWTDecode is the method to verify signature and decode JWT
// JWT 验签
func JWTValid(secret, jwt string) (bool, error) {

	ss, sign := strings.Split(jwt, "."), ""
	if len(ss) != 3 || ss[2] == sign {
		// 可能需要 sleep 200ms，用以应对时间差探测
		return false, nil
	}

	hD, err := base64.URLEncoding.DecodeString(ss[0])
	if err != nil {
		return false, err
	}

	h := &header{}
	if err = json.Unmarshal(hD, h); err != nil {
		return false, err
	}

	if h.Alg == "HS256" {
		sign, err = jwtHS256(secret, fmt.Sprintf("%s.%s", ss[0], ss[1]))
	} else if h.Alg == "RS256" {
		// sign = ras
	}

	if err != nil {
		return false, err
	}

	return sign == ss[2], nil
}

func JWTDecode(jwt string, m interface{}) error {
	ss := strings.Split(jwt, ".")
	if len(ss) < 3 {
		return errors.New("JWT Decode Error")
	}
	data, err := base64.URLEncoding.DecodeString(ss[1])
	if err != nil {
		return err
	}
	return json.Unmarshal(data, m)
}

// 签名方法
func jwtHS256(secret string, data string) (string, error) {

	s := sha256.New()
	_, err := s.Write([]byte(data))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum([]byte(secret))), nil
}

func jwtRS256() (string, error) {
	return "", nil
}
