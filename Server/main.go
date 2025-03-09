package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
)

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (Change "*" to specific frontend URL in production)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r) // Serve actual request
	})
}

var allowedExtensions = map[string]bool{
	".mp4": true,
}

// Handlers
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 100MB
	r.ParseMultipartForm(100 << 20) // 100MB

	// Get uploaded file - video (key name)
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Extract file extension
	ext := filepath.Ext(handler.Filename)
	if !allowedExtensions[ext] {
		http.Error(w, "Invalid file type. Only MP4 allowed.", http.StatusBadRequest)
		return
	}

	// Save uploaded file
	uploadPath := fmt.Sprintf("uploads/%s", handler.Filename)
	dst, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	// copy file to the path
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}

	// Generate output paths
	encryptedPath := "uploads/encrypted/"
	mpdPath := "uploads/encrypted/stream.mpd"

	// Run encryption
	err = EncryptDashAndPackage(uploadPath, encryptedPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded and encrypted successfully! MPD: %s", mpdPath)
}

func EncryptDashAndPackage(inputFile, outputFile string) error {
	cmd := exec.Command(
		"packager",
		fmt.Sprintf("input=%s,stream=video,output=%s/video.mp4", inputFile, outputFile),
		fmt.Sprintf("input=%s,stream=audio,output=%s/audio.mp4", inputFile, outputFile),
		"--mpd_output", fmt.Sprintf("%s/stream.mpd", outputFile),
		"--base_urls", "http://localhost:8080/encrypted/", // tells if the client is not on the same server whaere to pick the files
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Shaka Packager error: %v\n%s", err, string(output))
	}

	fmt.Println("Encryption & DASH Packaging completed successfully!")
	return nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to DRM Encoder")
}

func main() {
	r := mux.NewRouter()
	os.MkdirAll("uploads", os.ModePerm)
	// routes
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/upload", UploadHandler).Methods("POST")
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", enableCORS(http.FileServer(http.Dir("uploads/encrypted")))))

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
