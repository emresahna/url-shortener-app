package service

import (
	"context"
	"testing"

	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUrlRedirect(t *testing.T) {
	ctx := context.Background()
	code := "testcode"
	originalURL := "https://example.com"

	t.Run("successful_redirect_from_cache", func(t *testing.T) {
		mockDB := new(MockQuerier)
		mockRCC := new(MockRedisStore)
		mockRCA := new(MockRedisStore)
		mockAuth := new(MockAuth) // I'll need a mock for auth too

		s := New(mockDB, mockRCC, mockAuth, mockRCA)

		mockRCC.On("GetUrl", ctx, code).Return(originalURL, nil)
		mockRCA.On("IncreaseClick", ctx, code).Return(nil)

		res, err := s.UrlRedirect(ctx, code)

		assert.NoError(t, err)
		assert.Equal(t, originalURL, res)
		mockRCC.AssertExpectations(t)
		mockRCA.AssertExpectations(t)
	})

	t.Run("successful_redirect_from_db", func(t *testing.T) {
		mockDB := new(MockQuerier)
		mockRCC := new(MockRedisStore)
		mockRCA := new(MockRedisStore)
		mockAuth := new(MockAuth)

		s := New(mockDB, mockRCC, mockAuth, mockRCA)

		mockRCC.On("GetUrl", ctx, code).Return("", nil) // Cache miss
		mockDB.On("GetURLByCode", ctx, code).Return(originalURL, nil)
		mockRCA.On("IncreaseClick", ctx, code).Return(nil)

		res, err := s.UrlRedirect(ctx, code)

		assert.NoError(t, err)
		assert.Equal(t, originalURL, res)
		mockDB.AssertExpectations(t)
		mockRCC.AssertExpectations(t)
		mockRCA.AssertExpectations(t)
	})

	t.Run("url_not_found", func(t *testing.T) {
		mockDB := new(MockQuerier)
		mockRCC := new(MockRedisStore)
		mockRCA := new(MockRedisStore)
		mockAuth := new(MockAuth)

		s := New(mockDB, mockRCC, mockAuth, mockRCA)

		mockRCC.On("GetUrl", ctx, code).Return("", nil)
		mockDB.On("GetURLByCode", ctx, code).Return("", pgx.ErrNoRows)

		res, err := s.UrlRedirect(ctx, code)

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.Equal(t, models.UrlNotFoundErr().Error(), err.Error())
	})
}

// MockAuth is a mock of auth.Auth
type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) Create(user sqlc.User) (models.LoginUserResponse, error) {
	args := m.Called(user)
	return args.Get(0).(models.LoginUserResponse), args.Error(1)
}

func (m *MockAuth) Parse(token string) (jwt.MapClaims, error) {
	args := m.Called(token)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
