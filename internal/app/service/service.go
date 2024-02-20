package service

import (
	"context"

	"github.com/emarifer/gocms/internal/app/model"
	"github.com/emarifer/gocms/internal/app/repository"
)

// Service is the business logic of the application.
type Service interface {
	RecoverPosts(ctx context.Context) ([]model.Post, error)
	RecoverPost(ctx context.Context, id int) (*model.Post, error)
}

type serv struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {

	return &serv{
		repo: repo,
	}
}
