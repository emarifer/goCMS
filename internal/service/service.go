package service

import (
	"context"

	"github.com/emarifer/gocms/internal/entity"
	"github.com/emarifer/gocms/internal/model"
	"github.com/emarifer/gocms/internal/repository"
)

// Service is the business logic of the application.
type Service interface {
	CreatePost(ctx context.Context, post *entity.Post) (int, error)
	RecoverPosts(ctx context.Context) ([]model.Post, error)
	RecoverPost(ctx context.Context, id int) (*model.Post, error)
	ChangePost(ctx context.Context, post *entity.Post) (int64, error)
	RemovePost(ctx context.Context, id int) (int64, error)

	RecoverImageMetadata(ctx context.Context, uuid string) (*model.Image, error)
	CreateImageMetadata(ctx context.Context, image *entity.Image) error
	RemoveImage(ctx context.Context, uuid string) (int64, error)
}

type serv struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {

	return &serv{
		repo: repo,
	}
}
