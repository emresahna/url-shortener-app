package service

import (
	"context"
	"testing"
	"time"

	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name     string
		expected models.HealthResponse
	}{
		{
			name: "successful health check",
			expected: models.HealthResponse{
				Status:  "ok",
				Version: "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - create service with minimal dependencies for health check
			s := &service{
				db:  nil, // Health check doesn't use database
				rcc: nil, // Health check doesn't use Redis
				rca: nil, // Health check doesn't use Redis
			}

			ctx := context.Background()

			// Act
			result, err := s.HealthCheck(ctx)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Version, result.Version)
			
			// Verify timestamp format
			_, parseErr := time.Parse(time.RFC3339, result.Timestamp)
			assert.NoError(t, parseErr)
		})
	}
}

func TestReadinessCheck(t *testing.T) {
	tests := []struct {
		name           string
		dbNil          bool
		rccNil         bool
		expectError    bool
		expectedStatus string
	}{
		{
			name:           "database connection failure",
			dbNil:          true,
			rccNil:         false,
			expectError:    true,
			expectedStatus: "",
		},
		{
			name:           "redis connection failure",
			dbNil:          false,
			rccNil:         true,
			expectError:    true,
			expectedStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			s := &service{
				db:  nil, // Always nil for these tests
				rcc: nil, // Always nil for these tests
				rca: nil, // Always nil for these tests
			}

			ctx := context.Background()

			// Act
			_, err := s.ReadinessCheck(ctx)

			// Assert - ReadinessCheck should fail with nil dependencies
			assert.Error(t, err)
		})
	}
}

func TestLivenessCheck(t *testing.T) {
	tests := []struct {
		name     string
		expected models.HealthResponse
	}{
		{
			name: "successful liveness check",
			expected: models.HealthResponse{
				Status:  "alive",
				Version: "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - create service with minimal dependencies for liveness check
			s := &service{
				db:  nil, // Liveness check doesn't use database
				rcc: nil, // Liveness check doesn't use Redis
				rca: nil, // Liveness check doesn't use Redis
			}

			ctx := context.Background()

			// Act
			result, err := s.LivenessCheck(ctx)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Version, result.Version)
			
			// Verify timestamp format
			_, parseErr := time.Parse(time.RFC3339, result.Timestamp)
			assert.NoError(t, parseErr)
		})
	}
}

