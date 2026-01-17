package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"ris.com/internal/models"
)

const (
	// TODO: update to non-localhost targeting embedEndpoint
	embedEndpoint = "http://localhost:8000/api/v1/embed"
	expectedImageFieldName = "image"

	defaultEncodeQuality = 90
)

// EmbeddingGateway defines interaction with the image embedding service.
type EmbeddingGateway interface {
	GetEmbedding(ctx context.Context, image image.Image) (*models.EmbedResponse, error)
}

type embeddingGateway struct {
	client *http.Client
}

var _ EmbeddingGateway = &embeddingGateway{}

// NewEmbeddingGateway returns a new instance of the concrete implementation of EmbeddingGateway
func NewEmbeddingGateway(client *http.Client) EmbeddingGateway {
	return &embeddingGateway{client: client}
}

// GetEmbedding makes a POST request to the embedding service with the input image and returns a model containing the normalized vector embedding
func (eg *embeddingGateway) GetEmbedding(ctx context.Context, image image.Image) (*models.EmbedResponse, error) {
	// TODO: encode the image based on MIME type
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, image, &jpeg.Options{Quality: defaultEncodeQuality}); err != nil {
		return nil, err
	}

	// Open a new buffer with multipart.Writer to copy image contents and act as request body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(expectedImageFieldName, fmt.Sprintf("image:%s", uuid.NewString()))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, buf); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	// Create and execute the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", embedEndpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	resp, err := eg.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Server returned non-OK status: %d %s", resp.StatusCode, bodyBytes)
	}

	var embedding models.EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
		return nil, err
	}

	log.Printf("Embedding: %v", embedding.Embedding)
	return &embedding, nil
}

