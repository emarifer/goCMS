package repository

import (
	"context"

	"github.com/emarifer/gocms/internal/app/entity"
	"github.com/jmoiron/sqlx"
)

// Repository is the interface that wraps the basic CRUD operations.
type Repository interface {
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPost(ctx context.Context, id int) (*entity.Post, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {

	return &repo{
		db: db,
	}
}
