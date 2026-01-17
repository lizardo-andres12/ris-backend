package repository

import (
	"context"
	"database/sql"

	"github.com/pgvector/pgvector-go"
	"ris.com/internal/models"
)

const (
	similarityQuery = `SELECT id, domain_url, product_url, product_name, product_seller, product_price, file_name, 
file_size, height, width, format, 1 - (embedding <=> $1) AS cosine_similarity FROM images ORDER BY 
cosine_similarity DESC LIMIT $2 OFFSET $3`
)

// Interface for image table database operations
type ImageRepository interface {
	// QuerySimilar returns a slice of images similar to the input embedding ordered from most similar to least similar
	QuerySimilar(ctx context.Context, embedding []float32, limit, offset int) ([]*models.Image, error)
}

type imageRepository struct {
	db *sql.DB
}

var _ ImageRepository = &imageRepository{}

// NewImageRepository returns a new instance of the concrete implementation of ImageRepository
func NewImageRepository(db *sql.DB) ImageRepository {
	return &imageRepository{db: db}
}


func (ir *imageRepository) QuerySimilar(ctx context.Context, embedding []float32, limit, offset int) ([]*models.Image, error) {
	stmt, err := ir.db.PrepareContext(ctx, similarityQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	vector := pgvector.NewVector(embedding)
	rows, err := stmt.QueryContext(ctx, vector, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]*models.Image, limit, 0)
	for rows.Next() {
		var image models.Image
		if err := rows.Scan(&image.ID, &image.DomainURL, &image.ProductURL, &image.ProductName, &image.ProductSeller, &image.ProductPrice, &image.Filename, &image.Filesize, &image.Height, &image.Width, &image.Format, &image.Similarity); err != nil {
			return nil, err
		}
		images = append(images, &image)
	}

	return images, nil
}

