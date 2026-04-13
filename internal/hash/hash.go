package hash

import (
	"crypto/rand"
	"math/big"
)

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	codeLength  = 7
)

func toBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}
	result := ""
	for num > 0 {
		result = string(base62Chars[num%62]) + result
		num /= 62
	}
	return result
}

func GenerateUniqueCode() string {
	// 62^7 is approx 3.5 trillion
	max := new(big.Int).Exp(big.NewInt(62), big.NewInt(int64(codeLength)), nil)
	n, _ := rand.Int(rand.Reader, max)

	code := toBase62(n.Int64())

	// Pad with '0' if necessary to maintain consistent length
	for len(code) < codeLength {
		code = string(base62Chars[0]) + code
	}

	return code
}
