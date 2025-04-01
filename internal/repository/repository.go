package repository

import (
	"context"

	"github.com/spanwalla/url-shortener/internal/entity"
	"github.com/spanwalla/url-shortener/internal/repository/memory"
	pgRepo "github.com/spanwalla/url-shortener/internal/repository/postgres"
	"github.com/spanwalla/url-shortener/pkg/postgres"
)

//go:generate go tool mockgen -source=repository.go -destination=../mocks/repository/mock.go -package=repomocks

type Link interface {
	Store(ctx context.Context, alias, uri string) (string, error)
	Get(ctx context.Context, alias string) (entity.Link, error)
}

type Repositories struct {
	Link
}

func NewMemoryRepositories() *Repositories {
	return &Repositories{
		Link: memory.NewLinkRepo(),
	}
}

func NewPostgresRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Link: pgRepo.NewLinkRepo(pg),
	}
}
