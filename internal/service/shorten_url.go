package service

import (
	"context"
	"github.com/EmreSahna/url-shortener-app/internal/hash"
	"github.com/EmreSahna/url-shortener-app/internal/models"
	"github.com/EmreSahna/url-shortener-app/internal/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *service) ShortenURL(ctx context.Context, req models.ShortenURLRequest) (res models.ShortenURLResponse, err error) {
	t := ctx.Value("token").(string)
	c, err := s.jwt.Parse(t)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Parse user_id from token
	userUUID, err := uuid.Parse(c["id"].(string))
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Shorten url
	shortenUrl := hash.GenerateUniqueCode()

	// Save to db
	savedUrl, err := s.db.CreateURL(ctx, sqlc.CreateURLParams{
		OriginalUrl:   req.OriginalUrl,
		ShortenedCode: shortenUrl,
		UserID:        userUUID,
		ExpireTime: pgtype.Timestamptz{
			Time:  req.ExpireTime,
			Valid: true,
		},
	})
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Save to cache
	err = s.rc.SetUrl(ctx, shortenUrl, req.OriginalUrl)
	if err != nil {
		return models.ShortenURLResponse{}, err
	}

	// Return
	return models.ShortenURLResponse{
		Url: savedUrl.ShortenedCode,
	}, nil
}
