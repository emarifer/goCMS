package service

import (
	"context"

	"github.com/emarifer/gocms/internal/app/model"
)

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

func (s *serv) RecoverPost(
	ctx context.Context, id int,
) (*model.Post, error) {
	entityPost, err := s.repo.GetPost(ctx, id)
	if err != nil {
		/* if strings.Contains(err.Error(), "no rows in result set") {
			return nil, ErrResourceNotFound
		} */

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
