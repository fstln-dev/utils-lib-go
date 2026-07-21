package rand

import (
	"crypto/rand"
	"fmt"
)

func GetRandString(n int) string {
	randBytes := make([]byte, (n/2)+1)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)[:n+1]
}
