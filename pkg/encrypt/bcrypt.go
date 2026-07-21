package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// 加密
func Bcrypt(secret string, cost int) (string, error) {

	// bcrypt 加密
	res, err := bcrypt.GenerateFromPassword([]byte(secret), cost)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// 比对
func BcryptCompare(ori, input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(ori), []byte(input)) == nil
}
