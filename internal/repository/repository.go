package repository

import (
	"context"

	"github.com/emarifer/gocms/internal/entity"
	"github.com/jmoiron/sqlx"
)

// Repository is the interface that wraps the basic CRUD operations.
type Repository interface {
	SavePost(ctx context.Context, post *entity.Post) (int, error)
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPost(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post *entity.Post) (int64, error)
	DeletePost(ctx context.Context, id int) (int64, error)

	GetImage(ctx context.Context, uuid string) (*entity.Image, error)
	SaveImage(ctx context.Context, imageData *entity.Image) error
	DeleteImage(ctx context.Context, uuid string) (int64, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {

	return &repo{
		db: db,
	}
}
