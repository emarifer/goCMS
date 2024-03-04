package service

import (
	"context"
	"fmt"

	"github.com/emarifer/gocms/internal/entity"
	"github.com/emarifer/gocms/internal/model"
)

func (s *serv) RecoverImageMetadata(
	ctx context.Context, uuid string,
) (*model.Image, error) {
	entityImage, err := s.repo.GetImage(ctx, uuid)
	if err != nil {

		return nil, err
	}

	mi := &model.Image{
		UUID: entityImage.UUID.String(),
		Name: entityImage.Name,
		Alt:  entityImage.Alt,
	}

	return mi, nil
}

func (s *serv) RecoverAllImageMetadata(ctx context.Context, limit int) (
	[]model.Image, error,
) {
	ii := []model.Image{}

	entityImages, err := s.repo.GetImages(ctx, limit)
	if err != nil {
		return nil, err
	}

	for _, item := range entityImages {
		i := model.Image{
			UUID: item.UUID.String(),
			Name: item.Name,
			Alt:  item.Alt,
		}

		ii = append(ii, i)
	}

	return ii, nil
}

func (s *serv) CreateImageMetadata(
	ctx context.Context, image *entity.Image,
) error {
	if image.Name == "" {
		return fmt.Errorf("cannot have empty name")
	}

	if image.Alt == "" {
		return fmt.Errorf("cannot have empty alt text")
	}

	if err := s.repo.SaveImage(ctx, image); err != nil {
		return err
	}

	return nil
}

func (s *serv) RemoveImage(ctx context.Context, uuid string) (int64, error) {
	row, err := s.repo.DeleteImage(ctx, uuid)
	if err != nil {
		return row, err
	}

	return row, nil
}
