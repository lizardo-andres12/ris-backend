package models

import "github.com/pgvector/pgvector-go"

// Image represents an indexed product image with metadata saved at the latest pdate
type Image struct {
	ID int64 `db:"id" json:"id"`
	Embedding pgvector.Vector `db:"embedding"`
	Similarity float32 `db:"cosine_similarity" json:"similarity"`
	
	DomainURL string `db:"domain_url" json:"domain_url"`
	ProductURL string `db:"product_url" json:"product_url"`
	ProductName string `db:"product_name" json:"product_name"`
	ProductSeller string `db:"product_seller" json:"product_seller"`
	ProductPrice string `db:"product_price" json:"product_price"`

	Filename string `db:"file_name" json:"file_name"`
	Filesize int32 `db:"file_size" json:"file_size"`
	Height int32 `db:"height" json:"height"`
	Width int32 `db:"width" json:"width"`
	Format string `db:"format" json:"format"`
}

