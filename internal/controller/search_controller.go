package controller

import (
	"context"
	"image"

	"ris.com/internal/gateway"
	"ris.com/internal/models"
	"ris.com/internal/repository"
)

// SearchController manages all business logic and logging for similar image search functionality
type SearchController interface {
	// SearchSimilar requests a normalized embedding using its embeddingGate, queries similar images using that embedding,
	// and returns a list of images similar to the embedding
	SearchSimilar(ctx context.Context, image image.Image, limit, offset int) ([]*models.Image, error)
}

type searchController struct {
	imageRepository repository.ImageRepository
	embeddingGateway gateway.EmbeddingGateway
}

var _ SearchController = &searchController{}

// NewSearchController returns a new concrete implmenetation of the SearchController interface
func NewSearchController(ir repository.ImageRepository, eg gateway.EmbeddingGateway) SearchController {
	return &searchController{imageRepository: ir, embeddingGateway: eg}
}

func (sc *searchController) SearchSimilar(ctx context.Context, image image.Image, limit, offset int) ([]*models.Image, error) {
	embedding, err := sc.embeddingGateway.GetEmbedding(ctx, image)
	if err != nil {
		return nil, err
	}

	images, err := sc.imageRepository.QuerySimilar(ctx, embedding.Embedding, limit, offset)
	if err != nil {
		return nil, err
	}
	return images, nil
}

