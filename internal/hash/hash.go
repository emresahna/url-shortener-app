package hash

import (
	"github.com/google/uuid"
)

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func toBase62(num int64) string {
	return string(base62Chars[num%62])
}

func GenerateUniqueCode() string {
	rand := uuid.New().String()

	chunkSize := len(rand) / 7
	var result string
	for i := 0; i < 7; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(rand) {
			end = len(rand)
		}
		chunk := rand[start:end]

		var sum int64
		for _, c := range chunk {
			sum += int64(c)
		}

		result += toBase62(sum)
	}

	return result
}
