package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net/http"
	"time"

	"ris.com/internal/controller"
	"ris.com/internal/models"
)

const (
	maxFormSize = 10 << 20 // 10 MB

	imageFileKey = "image"
	contentTypeKey = "Content-Type"
	jsonEncondingHeader = "application/json"
)

// SearchSimilarHandler implements http.Handler and manages the searchSimilar endpoint, calling the SearchSimilar controller
type searchSimilarHandler struct {
	searchController controller.SearchController
}

func NewSearchSimilarHandler(sc controller.SearchController) http.Handler {
	return &searchSimilarHandler{searchController: sc}
}

func (ssh *searchSimilarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	file, fileHeaders, err := r.FormFile(imageFileKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find file for key %s: %v", imageFileKey, err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	image, format, err := image.Decode(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode image for given format %s: %v", format, err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	images, err := ssh.searchController.SearchSimilar(ctx, image, 0, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Server error: %v", err), http.StatusInternalServerError)
	}

	w.Header().Set(contentTypeKey, jsonEncondingHeader)
	response := models.SearchSimilarResponse {
		Images: images,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Server failed to encode response: %v", err), http.StatusInternalServerError)
	}
	
	log.Println(format)
	log.Printf("Content-Type: %s, Filename: %s, Filesize: %d", fileHeaders.Header.Get(contentTypeKey), fileHeaders.Filename, fileHeaders.Size)
}

