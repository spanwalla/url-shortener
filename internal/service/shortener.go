package service

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spanwalla/url-shortener/internal/repository"
	"github.com/spanwalla/url-shortener/internal/repository/repoerrs"
	"github.com/spanwalla/url-shortener/pkg/encoder"
)

type ShortenerService struct {
	linkRepo            repository.Link
	encoder             encoder.Encoder
	aliasLength         int
	attemptsOnCollision int
}

func NewShortenerService(linkRepo repository.Link, encoder encoder.Encoder, aliasLength, attemptsOnCollision int) *ShortenerService {
	return &ShortenerService{
		linkRepo:            linkRepo,
		encoder:             encoder,
		aliasLength:         aliasLength,
		attemptsOnCollision: attemptsOnCollision,
	}
}

func (s *ShortenerService) Shorten(ctx context.Context, uri string) (string, bool, error) {
	currAttempts := s.attemptsOnCollision

	for currAttempts > 0 {
		alias := s.encoder.Encode(uri, s.aliasLength)
		newAlias, err := s.linkRepo.Store(ctx, alias, uri)
		switch {
		case err == nil:
			return newAlias, alias == newAlias, nil
		case errors.Is(err, repoerrs.ErrAlreadyExists):
			log.Errorf("ShortenerService.Shorten - s.encoder.Encode: collision (%v, %v).", alias, uri)
			currAttempts--
			continue
		default:
			log.Errorf("ShortenerService.Shorten - s.encoder.Encode: %v", err)
			return "", false, ErrCannotShortenURI
		}
	}

	return "", false, ErrCannotShortenURI
}
