package hash

import (
	"math/rand"
	"time"
)

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

func toBase62(num int) string {
	return string(base62Chars[num%62])
}

func GenerateUniqueCode() string {
	chunkSize := 5
	var result string
	for i := 0; i < 7; i++ {
		var sum int
		for i := 0; i < chunkSize; i++ {
			sum += r.Intn(100)
		}

		result += toBase62(sum)
	}

	return result
}
