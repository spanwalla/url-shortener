package memory

import (
	"context"

	"github.com/spanwalla/url-shortener/internal/entity"
	"github.com/spanwalla/url-shortener/internal/repository/repoerrs"
	"github.com/spanwalla/url-shortener/pkg/memory"
)

type LinkRepo struct {
	uriToAlias *memory.Storage[string, string]
	aliasToURI *memory.Storage[string, string]
}

func NewLinkRepo() *LinkRepo {
	return &LinkRepo{
		uriToAlias: memory.NewStorage[string, string](), //
		aliasToURI: memory.NewStorage[string, string](),
	}
}

func (r *LinkRepo) Store(_ context.Context, alias, uri string) (string, error) {
	if val, ok := r.uriToAlias.Get(uri); ok {
		return val, nil
	}

	if _, ok := r.aliasToURI.Get(alias); ok {
		return "", repoerrs.ErrAlreadyExists
	}

	r.uriToAlias.Set(uri, alias)
	r.aliasToURI.Set(alias, uri)

	return alias, nil
}

func (r *LinkRepo) Get(_ context.Context, alias string) (entity.Link, error) {
	if val, ok := r.aliasToURI.Get(alias); ok {
		return entity.Link{Alias: alias, URI: val}, nil
	}

	return entity.Link{}, repoerrs.ErrNotFound
}
