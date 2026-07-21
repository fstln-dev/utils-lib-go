package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func HS256(secret string, data ...string) (string, error) {

	s := sha256.New()

	for _, d := range data {
		_, err := s.Write([]byte(d))
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(s.Sum([]byte(secret))), nil
}
