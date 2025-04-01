package service

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/spanwalla/url-shortener/internal/repository"
	"github.com/spanwalla/url-shortener/internal/repository/repoerrs"
)

type ExpanderService struct {
	linkRepo repository.Link
}

func NewExpanderService(linkRepo repository.Link) *ExpanderService {
	return &ExpanderService{
		linkRepo: linkRepo,
	}
}

func (s *ExpanderService) Expand(ctx context.Context, alias string) (string, error) {
	link, err := s.linkRepo.Get(ctx, alias)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrURINotFound
		}

		log.Errorf("ExpanderService.Expand - s.linkRepo.Get: %v", err)
		return "", ErrCannotExpandURI
	}

	return link.URI, nil
}
