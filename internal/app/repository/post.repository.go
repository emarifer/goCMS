package repository

import (
	"context"

	"github.com/emarifer/gocms/internal/app/entity"
)

const (
	qryGetPostById = `
		SELECT * FROM posts
		WHERE id = ?;
	`

	qryGetAllPosts = `
		SELECT * FROM posts
		ORDER BY created_at DESC;
	`
)

// This function gets all the posts from the current
// database
func (r *repo) GetPosts(ctx context.Context) ([]entity.Post, error) {
	pp := []entity.Post{}

	err := r.db.SelectContext(ctx, &pp, qryGetAllPosts)
	if err != nil {
		return nil, err
	}

	return pp, nil
}

// This function gets a post from the database
// with the given ID.
func (r *repo) GetPost(
	ctx context.Context, id int,
) (*entity.Post, error) {
	p := &entity.Post{}

	err := r.db.GetContext(ctx, p, qryGetPostById, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}
