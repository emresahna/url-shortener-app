package service

import (
	"context"
	"errors"
	"time"

	"github.com/emresahna/url-shortener-app/internal/hash"
	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/emresahna/url-shortener-app/internal/validator"
	"github.com/redis/go-redis/v9"
)

const (
	freeTierMinute   = 10
	freeTierDuration = time.Minute * freeTierMinute
)

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
