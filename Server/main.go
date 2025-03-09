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

var allowedEncryption = map[string]bool{
	"Widevine":  true,
	"RawKey":    true,
	"PlayReady": true,
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

	segmentSize := r.FormValue("segmentSize")
	encryptionType := r.FormValue("encryptionType")
	includeAudio := r.FormValue("includeAudio")

	if segmentSize == "" {
		segmentSize = "4" // Default segment size in seconds
	}
	if !allowedEncryption[encryptionType] {
		http.Error(w, "Invalid Encryption Type.", http.StatusBadRequest)
		return
	}

	// Generate output paths
	encryptedPath := "uploads/encrypted/"
	mpdPath := "uploads/encrypted/stream.mpd"

	// Run encryption
	err = EncryptDashAndPackage(uploadPath, encryptedPath, segmentSize, encryptionType, includeAudio)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded and encrypted successfully! MPD: %s", mpdPath)
}

func EncryptDashAndPackage(inputFile, outputFile, segmentSize, encryptionType, includeAudio string) error {

	cmdArgs := []string{
		"packager",
		fmt.Sprintf("input=%s,stream=video,init_segment=%s/video_init.mp4,segment_template=%s/video_$Number$.m4s", inputFile, outputFile, outputFile),
		fmt.Sprintf("--segment_duration=%s", segmentSize),
		"--segment_sap_aligned",
		"--generate_static_live_mpd",
		"--mpd_output", fmt.Sprintf("%s/stream.mpd", outputFile),
		"--base_urls", "http://localhost:8080/encrypted/",
	}

	if includeAudio == "yes" {
		cmdArgs = append(cmdArgs,
			fmt.Sprintf("input=%s,stream=audio,init_segment=%s/audio_init.mp4,segment_template=%s/audio_$Number$.m4s", inputFile, outputFile, outputFile),
		)
	}

	if encryptionType != "" {
		switch encryptionType {
		case "RawKey":
			cmdArgs = append(cmdArgs,
				"--enable_raw_key_encryption",
				"--keys", "key_id=07507c220e89a23e20b25a2d03b74d53:key=6e19d3fabf454e4f0be778844354cf81",
			)
		case "Widevine":
			cmdArgs = append(cmdArgs,
				"--enable_widevine_encryption",
				"--key_server_url=https://license.uat.widevine.com/cenc/getcontentkey/widevine_test",
				"--content_id=7465737420636f6e74656e74206964",
				"--signer=widevine_test",
				"--aes_signing_key=1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9",
				"--aes_signing_iv=d58ce954203b7c9a9a9d467f59839249",
				"--protection_systems=Widevine",
				"--keys", "key_id=07507c220e89a23e20b25a2d03b74d53:key=6e19d3fabf454e4f0be778844354cf81",
			)
		case "PlayReady":
			cmdArgs = append(cmdArgs,
				"--enable_raw_key_encryption",
				"--keys", "key_id=07507c220e89a23e20b25a2d03b74d53:key=6e19d3fabf454e4f0be778844354cf81",
			)
		default:
			log.Println("Unknown encryption type.")
		}
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
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
