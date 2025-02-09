package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		isValid bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Valid subdomain", "https://sub.example.com", true},
		{"Valid URL with port", "https://example.com:8080", true},
		{"Valid URL with path", "https://example.com/path", true},
		{"Valid URL with query params", "https://example.com/path?query=value", true},
		{"Valid URL with fragment", "https://example.com/path#section", true},
		{"Valid URL with dash in domain", "https://my-site.com", true},
		{"Valid URL with numbers in domain", "https://123example.com", true},
		{"Empty string", "", false},
		{"Missing scheme", "example.com", false},
		{"Invalid characters", "http://exa mple.com", false},
		{"Only scheme", "https://", false},
		{"Invalid TLD", "https://example.123", false},
		{"Space in URL", "https://ex ample.com", false},
		{"Multiple dots", "https://..example.com", false},
		{"Invalid protocol", "ftp://example.com", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid := ValidateURL(tc.url)
			assert.Equal(t, tc.isValid, valid)
		})
	}
}
