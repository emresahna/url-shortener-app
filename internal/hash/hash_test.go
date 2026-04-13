package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUniqueCode(t *testing.T) {
	code1 := GenerateUniqueCode()
	code2 := GenerateUniqueCode()

	assert.Len(t, code1, codeLength)
	assert.Len(t, code2, codeLength)
	assert.NotEqual(t, code1, code2)

	for _, char := range code1 {
		assert.Contains(t, base62Chars, string(char))
	}
}

func TestToBase62(t *testing.T) {
	assert.Equal(t, "0", toBase62(0))
	assert.Equal(t, "a", toBase62(36))
	assert.Equal(t, "z", toBase62(61))
}
