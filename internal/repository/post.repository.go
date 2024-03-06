package repository

import (
	"context"

	"github.com/emarifer/gocms/internal/entity"
)

const (
	qryInsertPost = `
		INSERT INTO posts (title, excerpt, content)
		VALUES (:title, :excerpt, :content);
	`

	qryGetPostById = `
		SELECT * FROM posts
		WHERE id = ?;
	`

	qryGetAllPosts = `
		SELECT * FROM posts
		ORDER BY created_at DESC;
	`

	qryUpdatePostTitle = `
		UPDATE posts
		SET title = ?
		WHERE id = ?;
	`

	qryUpdatePostExcerpt = `
		UPDATE posts
		SET excerpt = ?
		WHERE id = ?;
	`

	qryUpdatePostContent = `
		UPDATE posts
		SET content = ?
		WHERE id = ?;
	`

	qryDeletePost = `
		DELETE FROM posts
		WHERE id = ?;
	`
)

// SavePost inserts a record into the database
// passing the title, excerpt and content
func (r *repo) SavePost(ctx context.Context, post *entity.Post) (int, error) {
	result, err := r.db.NamedExecContext(ctx, qryInsertPost, post)
	if err != nil {
		return -1, err
	}

	lastId, err := result.LastInsertId() // SEE NOTE BELOW
	if err != nil {
		return -1, err
	}

	return int(lastId), nil
}

// GetPosts gets all the posts from the current
// database
func (r *repo) GetPosts(ctx context.Context) ([]entity.Post, error) {
	pp := []entity.Post{}

	err := r.db.SelectContext(ctx, &pp, qryGetAllPosts)
	if err != nil {
		return nil, err
	}

	return pp, nil
}

// GetPost gets a post from the database
// with the given ID.
func (r *repo) GetPost(ctx context.Context, id int) (*entity.Post, error) {
	p := &entity.Post{}

	err := r.db.GetContext(ctx, p, qryGetPostById, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// UpdatePost updates a post from the database.
// with the ID and data provided.
func (r *repo) UpdatePost(
	ctx context.Context, post *entity.Post,
) (int64, error) {
	var row int64
	tx, err := r.db.Beginx()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	if len(post.Title) > 0 {
		result, err := tx.ExecContext(
			ctx, qryUpdatePostTitle, post.Title, post.ID,
		)
		if err != nil {
			return -1, err
		}

		if row, err = result.RowsAffected(); err != nil {
			return -1, err
		} else if row == 0 {
			return -1, err
		}
	}

	if len(post.Excerpt) > 0 {
		result, err := tx.ExecContext(
			ctx, qryUpdatePostExcerpt, post.Excerpt, post.ID,
		)
		if err != nil {
			return -1, err
		}

		if row, err = result.RowsAffected(); err != nil {
			return -1, err
		} else if row == 0 {
			return -1, err
		}
	}

	if len(post.Content) > 0 {
		result, err := tx.ExecContext(
			ctx, qryUpdatePostContent, post.Content, post.ID,
		)
		if err != nil {
			return -1, err
		}

		if row, err = result.RowsAffected(); err != nil {
			return -1, err
		} else if row == 0 {
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return row, nil
}

// DeletePost delete a post from the database
// with the given ID.
func (r *repo) DeletePost(ctx context.Context, id int) (int64, error) {
	var row int64

	result, err := r.db.ExecContext(ctx, qryDeletePost, id)
	if err != nil {
		return -1, err
	}

	if row, err = result.RowsAffected(); err != nil {
		return -1, err
	} else if row == 0 {
		return -1, err
	}

	return row, nil
}

/* How to get id of last inserted row from sqlx?. SEE:
https://stackoverflow.com/questions/53990799/how-to-get-id-of-last-inserted-row-from-sqlx
*/
