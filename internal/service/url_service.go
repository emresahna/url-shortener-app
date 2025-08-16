package service

import (
	"context"
	"errors"
	"time"

	"github.com/emresahna/url-shortener-app/internal/hash"
	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/emresahna/url-shortener-app/internal/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/redis/go-redis/v9"
)

const (
	freeTierMinute   = 10
	freeTierDuration = time.Minute * freeTierMinute
)

// UrlRedirect redirects to the original URL based on the shortened code
func (s *service) UrlRedirect(ctx context.Context, code string) (res string, err error) {
	var original string

	// Check redis
	original, err = s.rcc.GetUrl(ctx, code)
	if !errors.Is(err, redis.Nil) && err != nil {
		return "", models.InternalServerErr()
	}

	// Redis not hit, check database
	if original == "" {
		original, err = s.db.GetURLByCode(ctx, code)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", models.UrlNotFoundErr()
			} else {
				return "", models.InternalServerErr()
			}
		}
	}

	// Increase click
	err = s.rca.IncreaseClick(ctx, code)
	if err != nil {
		return "", models.InternalServerErr()
	}

	// Return URL
	return original, nil
}

// UrlRemove removes a URL by its ID
func (s *service) UrlRemove(ctx context.Context, id string) (res models.RemoveUrlResponse, err error) {
	// Parse url id
	urlUUID, err := uuid.Parse(id)
	if err != nil {
		return models.RemoveUrlResponse{}, err
	}

	// Get url from db
	code, err := s.db.GetURLByID(ctx, urlUUID)
	if err != nil {
		return models.RemoveUrlResponse{}, err
	}

	// Delete from redis
	if err = s.rcc.DeleteUrl(ctx, code); err != nil {
		return models.RemoveUrlResponse{}, err
	}

	// Soft delete from db
	now := time.Now()
	if err = s.db.DeleteExpiredUrlByShortCode(ctx, sqlc.DeleteExpiredUrlByShortCodeParams{
		DeletedAt:     &now,
		ShortenedCode: code,
	}); err != nil {
		return models.RemoveUrlResponse{}, err
	}

	return models.RemoveUrlResponse{
		Message: "Url successfully removed.",
	}, nil
}

// UrlShortenGuest shortens URL for guest users (with free tier limitations)
func (s *service) UrlShortenGuest(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Validate url
	if !validator.ValidateURL(req.OriginalUrl) {
		return models.ShortenURLResponse{}, models.UrlNotValidErr()
	}

	// Get ip
	ipAddr, ok := ctx.Value(models.IpKey).(string)
	if !ok {
		return models.ShortenURLResponse{}, err
	}

	// Check ip addr to not exceed free tier
	isFirstUsage := false
	count, err := s.rcc.GetIpAddrUsage(ctx, ipAddr)
	if errors.Is(err, redis.Nil) {
		isFirstUsage = true
	} else if err != nil {
		return models.ShortenURLResponse{}, err
	}

	if count >= 3 {
		return models.ShortenURLResponse{}, models.FreeTierExceedErr()
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	// Save to redis
	redisCh := make(chan error, 1)
	go func() {
		err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, freeTierDuration)
		if err != nil {
			redisCh <- models.SaveToCacheErr()
		}
		redisCh <- err
	}()

	// Save to db
	dbCh := make(chan error, 1)
	go func() {
		savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
			OriginalUrl:   req.OriginalUrl,
			ShortenedCode: shortenUrl,
			IpAddress:     &ipAddr,
		})
		if err != nil {
			dbCh <- models.CreateURLErr()
		}

		err = s.db.InsertClickCount(context.Background(), savedUrl.ID)
		if err != nil {
			dbCh <- models.CreateClickCountErr()
		}
		dbCh <- nil
	}()

	// Increase free usage
	usageCh := make(chan error, 1)
	go func() {
		if isFirstUsage {
			err = s.rcc.SetIpAddrUsage(ctx, ipAddr)
			if err != nil {
				usageCh <- models.SaveToCacheErr()
			}
			usageCh <- err
		} else {
			err = s.rcc.IncreaseIpAddrUsage(ctx, ipAddr)
			if err != nil {
				usageCh <- models.SaveToCacheErr()
			}
			usageCh <- err
		}
	}()

	if err = <-dbCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	if err = <-redisCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	if err = <-usageCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: shortenUrl,
	}, nil
}

// UrlShortenUser shortens URL for authenticated users
func (s *service) UrlShortenUser(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Validate url
	if !validator.ValidateURL(req.OriginalUrl) {
		return models.ShortenURLResponse{}, models.UrlNotValidErr()
	}

	// Parse token from string
	t := ctx.Value(models.TokenKey).(string)
	c, err := s.jwt.Parse(t)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return models.ShortenURLResponse{}, models.TokenExpiredErr()
		} else if err != nil {
			return models.ShortenURLResponse{}, models.TokenFailureErr()
		}
	}

	// Parse user id from token
	userUUID, err := uuid.Parse(c["id"].(string))
	if err != nil {
		return models.ShortenURLResponse{}, models.CustomerIdParseErr()
	}

	// Validate duration
	var duration time.Duration
	if req.ExpireTime != "" {
		parsedDate, err := validator.ParseDateWithTimeZone(req.ExpireTime)
		if err != nil {
			return models.ShortenURLResponse{}, err
		}

		duration, err = validator.ValidateFutureDate(parsedDate)
		if err != nil {
			return models.ShortenURLResponse{}, err
		}
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	// Save to redis
	redisCh := make(chan error, 1)
	go func() {
		err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, duration)
		if err != nil {
			redisCh <- models.SaveToCacheErr()
		}
		redisCh <- err
	}()

	// Save to db
	dbCh := make(chan error, 1)
	go func() {
		savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
			OriginalUrl:   req.OriginalUrl,
			ShortenedCode: shortenUrl,
			UserID:        &userUUID,
		})
		if err != nil {
			dbCh <- models.CreateURLErr()
		}

		err = s.db.InsertClickCount(context.Background(), savedUrl.ID)
		if err != nil {
			dbCh <- models.CreateClickCountErr()
		}
		dbCh <- nil
	}()

	if err = <-dbCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	if err = <-redisCh; err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: shortenUrl,
	}, nil
}
