package service

import (
	"context"

	"github.com/spanwalla/url-shortener/internal/repository"
	"github.com/spanwalla/url-shortener/pkg/encoder"
)

//go:generate mockgen -source=service.go -destination=../mocks/service/mock.go -package=servicemocks

type Shortener interface {
	Shorten(ctx context.Context, uri string) (string, bool, error)
}

type Expander interface {
	Expand(ctx context.Context, alias string) (string, error)
}

type Services struct {
	Shortener
	Expander
}

type Dependencies struct {
	Repos               *repository.Repositories
	Encoder             encoder.Encoder
	AliasLength         int
	AttemptsOnCollision int
}

func New(deps Dependencies) *Services {
	return &Services{
		Shortener: NewShortenerService(deps.Repos.Link, deps.Encoder, deps.AliasLength, deps.AttemptsOnCollision),
		Expander:  NewExpanderService(deps.Repos.Link),
	}
}
