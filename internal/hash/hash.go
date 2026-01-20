package hash

import (
	"crypto/rand"
	"math/big"
	"time"
)

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	chunkSize   = 5
	codeLength  = 7
)

func toBase62(num int) string {
	return string(base62Chars[num%62])
}

func GenerateUniqueCode() string {
	result := ""
	for range codeLength {
		sum := 0
		for range chunkSize {
			r, _ := rand.Int(rand.Reader, big.NewInt(time.Now().Unix()))
			sum += int(r.Int64())
		}
		result += toBase62(sum)
	}

	return result
}
