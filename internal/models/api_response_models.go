package models

// SearchSimilarResponse nests multiple Image models for json encoding convenience
type SearchSimilarResponse struct {
	Images []*Image `json:"images"`
}

