package service

import (
	"context"

	"github.com/emarifer/gocms/internal/entity"
	"github.com/emarifer/gocms/internal/model"
)

func (s *serv) CreatePost(ctx context.Context, post *entity.Post) (int, error) {
	var (
		lastInsertId int
		err          error
	)

	if lastInsertId, err = s.repo.SavePost(ctx, post); err != nil {
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (s *serv) RecoverPosts(ctx context.Context) ([]model.Post, error) {
	pp := []model.Post{}

	entityPosts, err := s.repo.GetPosts(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range entityPosts {
		p := model.Post{
			ID:      item.ID,
			Title:   item.Title,
			Excerpt: item.Excerpt,
		}

		pp = append(pp, p)
	}

	return pp, nil
}

func (s *serv) RecoverPost(ctx context.Context, id int) (*model.Post, error) {
	entityPost, err := s.repo.GetPost(ctx, id)
	if err != nil {

		return nil, err
	}

	mp := &model.Post{
		ID:        entityPost.ID,
		Title:     entityPost.Title,
		Content:   entityPost.Content,
		Excerpt:   entityPost.Excerpt,
		CreatedAt: entityPost.CreatedAt,
	}

	return mp, nil
}

func (s *serv) ChangePost(
	ctx context.Context, post *entity.Post,
) (int64, error) {
	row, err := s.repo.UpdatePost(ctx, post)
	if err != nil {
		return row, err
	}

	return row, nil
}

func (s *serv) RemovePost(ctx context.Context, id int) (int64, error) {
	row, err := s.repo.DeletePost(ctx, id)
	if err != nil {
		return row, err
	}

	return row, nil
}
