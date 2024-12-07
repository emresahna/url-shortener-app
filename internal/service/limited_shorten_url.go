package service

import (
	"context"
	"errors"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/EmreSahna/url-shortener-app/internal/validator"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func (s *service) LimitedShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	// Get ip
	ipAddr, ok := ctx.Value("ip").(string)
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

	// Validate url
	if !validator.ValidateURL(req.OriginalUrl) {
		return models.ShortenURLResponse{}, models.UrlNotValidErr()
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	duration := time.Minute * 10

	// Save to redis
	redisCh := make(chan error, 1)
	go func() {
		log.Printf("Starting to save limited shorten URL %s to Redis...", shortenUrl)
		err = s.rcc.SetUrlWithExpire(ctx, shortenUrl, req.OriginalUrl, duration)
		if err != nil {
			redisCh <- models.SaveToCacheErr()
		}
		redisCh <- err
	}()

	// Save to db
	dbCh := make(chan error, 1)
	go func() {
		log.Printf("Starting to save limited shorten URL %s to PostgreSQL...", shortenUrl)
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
		log.Printf("Starting to increase free usage for %s...", ipAddr)
		if isFirstUsage {
			err = s.rcc.IncreaseIpAddrUsage(ctx, ipAddr)
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
	log.Printf("Successfully saved limited shorten URL %s to Redis.", shortenUrl)

	if err = <-redisCh; err != nil {
		return models.ShortenURLResponse{}, err
	}
	log.Printf("Successfully saved limited shorten URL %s to PostgreSQL.", shortenUrl)

	if err = <-usageCh; err != nil {
		return models.ShortenURLResponse{}, err
	}
	log.Printf("Successfully increased free usage for %s.", ipAddr)

	// Return
	return models.ShortenURLResponse{
		Url: shortenUrl,
	}, nil
}
