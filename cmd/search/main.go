package main

import (
	"context"
	"log"
	"time"

	"ris.com/internal"
	"ris.com/internal/controller"
	"ris.com/internal/repository"
)

func main() {
	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	ir := repository.NewImageRepository(db)
	sc := controller.NewSearchController(ir)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second * 3)
	defer cancel()

	sc.SearchSimilar(ctx, nil, 5, 0)
}

func x() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. Parse the multipart form (allow up to 10MB)
		// This is required when the client sends data via requests.post(files={...})
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// 2. Retrieve the file. "image" matches the key in Python: files={'image': ...}
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to retrieve image file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 3. Decode the image directly from the form file
		img, format, err := image.Decode(file)
		if err != nil {
			log.Printf("Decode error: %v", err)
			http.Error(w, "Failed to decode image data", http.StatusBadRequest)
			return
		}

		// 4. Save the decoded image
		outFile, err := os.Create("decoded.jpg")
		if err != nil {
			http.Error(w, "Failed to create local file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
		if err != nil {
			http.Error(w, "Failed to save JPEG", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully received and saved image (format: %s)", format)
	})

	log.Println("Server starting on :65000...")
	if err := http.ListenAndServe(":65000", nil); err != nil {
		log.Fatal("Failed to start server")
	}
}


