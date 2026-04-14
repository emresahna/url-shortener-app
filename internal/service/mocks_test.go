package service

import (
	"context"
	"time"

	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock of sqlc.Querier
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateURL(ctx context.Context, arg sqlc.CreateURLParams) (sqlc.CreateURLRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.CreateURLRow), args.Error(1)
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) DeleteExpiredUrlByShortCode(ctx context.Context, arg sqlc.DeleteExpiredUrlByShortCodeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetIDByShortCode(ctx context.Context, shortenedCode string) (uuid.UUID, error) {
	args := m.Called(ctx, shortenedCode)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockQuerier) GetURLByCode(ctx context.Context, shortenedCode string) (string, error) {
	args := m.Called(ctx, shortenedCode)
	return args.String(0), args.Error(1)
}

func (m *MockQuerier) GetURLByID(ctx context.Context, id uuid.UUID) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockQuerier) GetUrlsByUser(ctx context.Context, arg sqlc.GetUrlsByUserParams) ([]sqlc.GetUrlsByUserRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetUrlsByUserRow), args.Error(1)
}

func (m *MockQuerier) GetUserByUserID(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) GetUserByUsername(ctx context.Context, username string) (sqlc.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) IncrementClickCount(ctx context.Context, arg sqlc.IncrementClickCountParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) InsertClickCount(ctx context.Context, urlID uuid.UUID) error {
	args := m.Called(ctx, urlID)
	return args.Error(0)
}

func (m *MockQuerier) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

// MockRedisStore is a mock of redis.Store
type MockRedisStore struct {
	mock.Mock
}

func (m *MockRedisStore) GetUrl(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

func (m *MockRedisStore) SetUrlWithExpire(ctx context.Context, code string, url string, expiration time.Duration) error {
	args := m.Called(ctx, code, url, expiration)
	return args.Error(0)
}

func (m *MockRedisStore) DeleteUrl(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *MockRedisStore) GetIpAddrUsage(ctx context.Context, ip string) (int, error) {
	args := m.Called(ctx, ip)
	return args.Int(0), args.Error(1)
}

func (m *MockRedisStore) SetIpAddrUsage(ctx context.Context, ip string) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}

func (m *MockRedisStore) IncreaseIpAddrUsage(ctx context.Context, ip string) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}

func (m *MockRedisStore) IncreaseClick(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *MockRedisStore) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
