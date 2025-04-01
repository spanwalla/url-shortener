package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/spanwalla/url-shortener/internal/entity"
	"github.com/spanwalla/url-shortener/internal/repository/repoerrs"
	"github.com/spanwalla/url-shortener/pkg/postgres"
)

type LinkRepo struct {
	*postgres.Postgres
}

func NewLinkRepo(pg *postgres.Postgres) *LinkRepo {
	return &LinkRepo{pg}
}

func (r *LinkRepo) Store(ctx context.Context, alias, uri string) (string, error) {
	sql, args, _ := r.Builder.
		Insert("links").
		Columns("alias, uri").
		Values(alias, uri).
		Suffix("ON CONFLICT (uri) DO UPDATE SET uri = EXCLUDED.uri"). // По сути присваивается то же самое значение,
		// но из-за этого мы можем получить alias через RETURNING
		Suffix("RETURNING alias").
		ToSql()

	var newAlias string
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&newAlias)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return "", repoerrs.ErrAlreadyExists
			}
		}
		return "", fmt.Errorf("LinkRepo.Store - r.Pool.QueryRow: %w", err)
	}

	return newAlias, nil
}

func (r *LinkRepo) Get(ctx context.Context, alias string) (entity.Link, error) {
	sql, args, _ := r.Builder.
		Select("uri").
		From("links").
		Where("alias = ?", alias).
		Limit(1).
		ToSql()

	link := entity.Link{Alias: alias}
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&link.URI)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Link{}, repoerrs.ErrNotFound
		}
		return entity.Link{}, fmt.Errorf("LinkRepo.Get - r.Pool.QueryRow: %w", err)
	}

	return link, nil
}
