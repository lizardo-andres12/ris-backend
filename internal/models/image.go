package models

import "github.com/pgvector/pgvector-go"

// Image represents an indexed product image with metadata saved at the latest pdate
type Image struct {
	ID int64 `db:"id"`
	Embedding pgvector.Vector `db:"embedding"`
	
	DomainURL string `db:"domain_url"`
	ProductURL string `db:"product_url"`
	ProductName string `db:"product_name"`
	ProductSeller string `db:"product_seller"`
	ProductPrice string `db:"product_price"`

	Filename string `db:"file_name"`
	Filesize int32 `db:"file_size"`
	Height int32 `db:"height"`
	Width int32 `db:"width"`
	Format string `db:"format"`
}

