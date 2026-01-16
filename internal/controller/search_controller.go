package controller

import (
	"context"
	"image"

	"ris.com/internal/models"
	"ris.com/internal/repository"
)

// SearchController manages all business logic and handles logging
type SearchController interface {
	SearchSimilar(ctx context.Context, image image.Image, limit, offset int) ([]*models.Image, error)
}

type searchController struct {
	imageRepository repository.ImageRepository
}

var _ SearchController = &searchController{}

func NewSearchController(ir repository.ImageRepository) SearchController {
	return &searchController{imageRepository: ir}
}

func (sc *searchController) SearchSimilar(ctx context.Context, image image.Image, limit, offset int) ([]*models.Image, error) {
//embedding, err := sc.embedderGateway.RequestEmbed(image)
//if err != nil {
//	// error log
//	return nil, err
//}
	embedding := []float32{1, 2, 3}
	images, err := sc.imageRepository.QuerySimilar(ctx, embedding, limit, offset)
	if err != nil {
		return nil, err
	}
	return images, nil
}

