package validator

import (
	"testing"
	"time"

	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/stretchr/testify/assert"
)

var mockNow = time.Date(2025, time.January, 0, 0, 0, 0, 0, time.UTC)
var mockFuture = time.Date(2699, time.January, 0, 0, 0, 0, 0, time.UTC)
var mockEarly = time.Date(1453, time.January, 0, 0, 0, 0, 0, time.UTC)

func TestParseDateWithTimeZone(t *testing.T) {
	tests := []struct {
		name           string
		date           string
		timeDifference time.Time
		expectedError  error
	}{
		{
			name:           "Time is valid (Now)",
			date:           mockNow.Format(time.RFC3339),
			timeDifference: mockNow,
			expectedError:  nil,
		},
		{
			name:           "Time is valid (Future)",
			date:           mockFuture.Format(time.RFC3339),
			timeDifference: mockFuture,
			expectedError:  nil,
		},
		{
			name:           "Time is valid (Past)",
			date:           mockEarly.Format(time.RFC3339),
			timeDifference: mockEarly,
			expectedError:  nil,
		},
		{
			name:          "Time is broken",
			date:          "0120-45-12-G35:12:444+0021",
			expectedError: models.ParseExpireTimeErr(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedDate, err := ParseDateWithTimeZone(tc.date)

			// assert that errors are correct
			assert.Equal(t, err, tc.expectedError)

			if tc.expectedError != nil {
				assert.Zero(t, parsedDate)
			} else {
				assert.Equal(t, parsedDate, tc.timeDifference)
			}
		})
	}
}

func TestValidateFutureDate(t *testing.T) {
	tests := []struct {
		name          string
		time          time.Time
		expectedError error
	}{
		{
			name:          "Time is valid (Future)",
			time:          mockFuture,
			expectedError: nil,
		},
		{
			name:          "Time is valid (Past)",
			time:          mockEarly,
			expectedError: models.ExpireTimeAlreadyPassedErr(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			duration, err := ValidateFutureDate(tc.time)

			// assert that errors are correct
			assert.Equal(t, err, tc.expectedError)

			if tc.expectedError != nil {
				assert.Zero(t, duration)
			} else {
				assert.Greater(t, duration, time.Duration(0))
			}
		})
	}
}
