package repository

import (
	"context"

	"github.com/emarifer/gocms/internal/entity"
)

const (
	qryInsertImageData = `
		INSERT INTO images (uuid, name, alt)
		VALUES (:uuid, :name, :alt);
	`

	qryGetImageById = `
		SELECT * FROM images
		WHERE uuid = ?;
	`

	qryDeleteImage = `
		DELETE FROM images
		WHERE uuid = ?;
	`
)

// This function inserts the metadata of an image
// in the database, passing the id, name and alt text
func (r *repo) SaveImage(ctx context.Context, imageData *entity.Image) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExecContext(ctx, qryInsertImageData, imageData)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// This function gets image metadata from the database
// with the given ID.
func (r *repo) GetImage(ctx context.Context, uuid string) (*entity.Image, error) {
	img := &entity.Image{}

	err := r.db.GetContext(ctx, img, qryGetImageById, uuid)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// This function delete image metadata from the database
// with the given ID.
func (r *repo) DeleteImage(ctx context.Context, uuid string) (int64, error) {
	var row int64

	result, err := r.db.ExecContext(ctx, qryDeleteImage, uuid)
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
