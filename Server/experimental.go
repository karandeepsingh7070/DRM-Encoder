package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath" // Needed for file extension checking
	"time"

	"github.com/gorilla/mux"
)

// var allowedExtensions = map[string]bool{
// 	".mp4": true,
// 	".mov": true,
// }

// CORS Middleware
func enableCORS2(h http.Handler) http.Handler {
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

// UploadHandler - Handles video uploads
// func UploadHandler(w http.ResponseWriter, r *http.Request) {
// 	// Limit upload size to 100MB
// 	r.ParseMultipartForm(100 << 20) // 100MB

// 	// Get the uploaded file
// 	file, handler, err := r.FormFile("video")
// 	if err != nil {
// 		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Extract file extension and check if it's allowed
// 	ext := filepath.Ext(handler.Filename)
// 	if !allowedExtensions[ext] {
// 		http.Error(w, "Invalid file type. Only MP4 and MOV are allowed.", http.StatusBadRequest)
// 		return
// 	}

// 	// Create a new file in "uploads/" directory
// 	// dst, err := os.Create(fmt.Sprintf("uploads/%s", handler.Filename))
// 	// if err != nil {
// 	// 	http.Error(w, "Error saving the file", http.StatusInternalServerError)
// 	// 	return
// 	// }
// 	// defer dst.Close()

// 	// Save original file
// 	originalPath := fmt.Sprintf("uploads/%s", handler.Filename)
// 	dst, err := os.Create(originalPath)
// 	if err != nil {
// 		http.Error(w, "Error saving the file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer dst.Close()

// 	// Copy uploaded file to new location
// 	_, err = io.Copy(dst, file)
// 	if err != nil {
// 		http.Error(w, "Error writing the file", http.StatusInternalServerError)
// 		return
// 	}
// 	// Create a progress channel
// 	// progressChan := make(chan string)
// 	// Convert video format
// 	outputPath := fmt.Sprintf("uploads/converted_h264_%s", handler.Filename) // Convert to MKV
// 	// err = ConvertVideo(originalPath, outputPath, "uploads", progressChan)
// 	ok := ConvertToMP4(originalPath, outputPath)
// 	fmt.Fprintln(w, "conversion started")
// 	if ok != nil {
// 		http.Error(w, "Video conversion failed", http.StatusInternalServerError)
// 		// close(progressChan)
// 		return
// 	}
// 	// _, err = io.Copy(dst, ok)

// 	fmt.Fprintf(w, "File uploaded and converted successfully: %s", outputPath)
// }

// ConvertVideo runs FFmpeg to convert a video to MKV
// func ConvertVideo(inputPath, outputPath string) error {
// 	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:v", "libx264", "-preset", "fast", outputPath) // equivalent to "ffmpeg -i input.mp4 -c:v libx264 -preset fast -crf 23 -c:a aac -b:a 128k output_h264.mp4"
// 	err := cmd.Run()
// 	return err
// }

// UploadHandler handles file uploads
func UploadHandler2(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 100MB
	r.ParseMultipartForm(100 << 20) // 100MB

	// Get uploaded file
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Extract file extension
	ext := filepath.Ext(handler.Filename)
	if ext != ".mp4" {
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
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}

	// Generate output paths
	encryptedPath := "uploads/encrypted/"
	mpdPath := "uploads/encrypted/stream.mpd"

	// Run encryption
	err = EncryptAndPackage(uploadPath, encryptedPath, mpdPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encryption failed: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded and encrypted successfully! MPD: %s", mpdPath)
}

func EncryptAndPackage(inputFile, outputFile, mpdFile string) error {
	// cmd := exec.Command("packager",
	// 	// fmt.Sprintf("in=%s,stream=video,output=%s,drm_label=SD", inputFile, outputFile),
	// 	fmt.Sprintf("in=%s,stream=video,init_segment=h264/init.mp4,segment_template=h264_segments/$Number$.m4s", inputFile),
	// 	// in=h264_baseline_360p_600.mp4,stream=video,init_segment=h264_360p/init.mp4,segment_template=h264_360p/$Number$.m4s
	// 	"--enable_raw_key_encryption",
	// 	"--keys=label=SD:key_id=0123456789abcdef0123456789abcdef:key=abcdef0123456789abcdef0123456789",
	// 	"--generate_static_live_mpd",
	// 	fmt.Sprintf("--mpd_output=%s", mpdFile),
	// )

	// cmd := exec.Command("packager",
	// 	fmt.Sprintf("in=%s,stream=video,init_segment=%s/video_init.mp4,segment_template=%s/video_segment_$Number$.m4s", inputFile, outputFile, outputFile),
	// 	fmt.Sprintf("in=%s,stream=audio,init_segment=%s/audio_init.mp4,segment_template=%s/audio_segment_$Number$.m4s", inputFile, outputFile, outputFile),
	// 	"--enable_raw_key_encryption",
	// 	"--keys", "key_id=1234567890abcdef1234567890abcdef:key=abcdef1234567890abcdef1234567890",
	// 	"--protection_scheme", "cbcs",
	// 	"--mpd_output", fmt.Sprintf("%s/stream.mpd", outputFile),
	// )

	cmd := exec.Command(
		"packager",
		fmt.Sprintf("input=%s,stream=video,output=%s/video.mp4", inputFile, outputFile),
		fmt.Sprintf("input=%s,stream=audio,output=%s/audio.mp4", inputFile, outputFile),
		"--mpd_output", fmt.Sprintf("%s/stream.mpd", outputFile),
		"--base_urls", "http://localhost:8080/encrypted/",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Shaka Packager error: %v\n%s", err, string(output))
	}

	fmt.Println("Encryption & DASH Packaging completed successfully!")
	return nil
}

// func PackageToDASH(inputFile, uploadFolder string) error {
// 	outputMPD := fmt.Sprintf("%s/output.mpd", uploadFolder)

// 	drmConfig := fmt.Sprintf("%s/drm.xml", "Server")
// 	// Debug: Print input file path
// 	fmt.Println("DASH Packaging Input File:", inputFile)
// 	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
// 		return fmt.Errorf("MP4Box error: input file does not exist: %s", inputFile)
// 	}

// 	// Run MP4Box
// 	cmd := exec.Command("MP4Box", "-crypt", drmConfig, "-dash", "4000", "-frag", "4000", "-rap",
// 		"-segment-name", fmt.Sprintf("%s/segment_", "segments"), // Store segments in uploads folder
// 		"-out", outputMPD,
// 		"-bs-switching", "no",
// 		"-profile", "dashavc264:onDemand",
// 		"-mpd-title", "DRM Video",
// 		"-subsegs-per-sidx", "1",
// 		"-crypt", drmConfig,
// 		inputFile)
// 	// cmd := exec.Command("MP4Box",
// 	// 	"-dash", "4000", "-frag", "4000", "-rap",
// 	// 	"-segment-name", fmt.Sprintf("%s/segment_", "segments"),
// 	// 	"-out", outputMPD,
// 	// 	"-encrypt",
// 	// 	"-enc-key", fmt.Sprintf("%s:%s", kid, key), // KID:KEY pair
// 	// 	inputFile,
// 	// )
// 	err := cmd.Run()
// 	return err

// 	// cmd.Dir = uploadFolder

// 	// output, err := cmd.CombinedOutput()
// 	// if err != nil {
// 	// 	return fmt.Errorf("MP4Box error: %v\n%s", err, string(output))
// 	// }

// 	// fmt.Println("DASH Packaging completed in", uploadFolder)
// 	// return nil
// }

// func ConvertToMP4(inputFile, outputFile string) error {
// 	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c:v", "libx264", "-c:a", "aac",
// 		"-movflags", "faststart+frag_keyframe", // Ensures MP4 is compatible with MP4Box
// 		outputFile)

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("FFmpeg error: %v\n%s", err, string(output))
// 	}

// 	fmt.Println("Video converted successfully! Checking if the file exists...")

// 	// Check if file exists before passing to MP4Box
// 	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
// 		return fmt.Errorf("Converted file does not exist: %s", outputFile)
// 	}

// 	// Now package it to DASH
// 	err = PackageToDASH(outputFile, "uploads")
// 	if err != nil {
// 		fmt.Println("Error while converting to DASH:", err)
// 	}
// 	return nil
// }

// func ConvertVideo(inputFile, outputFile, uploadFolder string, progressChan chan string) error {
// 	// Convert video using FFmpeg
// 	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", outputFile)

// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}

// 	if err := cmd.Start(); err != nil {
// 		return err
// 	}

// 	// Read FFmpeg progress and send updates
// 	scanner := bufio.NewScanner(stdout)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.Contains(line, "out_time_ms") {
// 			parts := strings.Split(line, "=")
// 			if len(parts) > 1 {
// 				progressChan <- fmt.Sprintf("FFmpeg Progress: %s", parts[1])
// 			}
// 		}
// 	}

// 	if err := cmd.Wait(); err != nil {
// 		return err
// 	}

// 	progressChan <- "Video conversion completed!"
// 	fmt.Println("Video converted successfully!")

// 	// **NEW: Ensure MP4Box stores segments inside the upload folder**
// 	dashManifest := fmt.Sprintf("%s/output.mpd", uploadFolder)
// 	dashCmd := exec.Command("MP4Box", "-dash", "4000", "-frag", "4000", "-rap",
// 		"-segment-name", fmt.Sprintf("%s/segment_", uploadFolder),
// 		"-out", dashManifest, outputFile)

// 	// Set the working directory to the upload folder
// 	dashCmd.Dir = uploadFolder

// 	dashStdout, err := dashCmd.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}

// 	if err := dashCmd.Start(); err != nil {
// 		return err
// 	}

// 	// Read MP4Box progress and send updates
// 	dashScanner := bufio.NewScanner(dashStdout)
// 	for dashScanner.Scan() {
// 		progressChan <- fmt.Sprintf("MP4Box Progress: %s", dashScanner.Text())
// 	}

// 	if err := dashCmd.Wait(); err != nil {
// 		return fmt.Errorf("MP4Box error: %v", err)
// 	}

// 	progressChan <- "DASH Packaging completed!"
// 	fmt.Println("DASH Packaging completed!")

// 	return nil
// }

func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	progressChan := make(chan string)

	// go ConvertVideo("uploads/sample.mp4", "uploads/sample.mkv", "uploads", progressChan)

	for progress := range progressChan {
		fmt.Fprintf(w, "data: %s\n\n", progress)
		flusher.Flush() // Send data to client
	}
}

// HomeHandler - Handles requests to "/"
// func HomeHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Welcome to the DRM Converter!")
// }

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get path variables
	videoID := vars["id"]
	fmt.Fprintf(w, "Video ID: %s", videoID)
}

// ProtectedHandler requires authentication
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You are authenticated!")
}

// AuthMiddleware checks if the request contains a valid API key
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key") // Get API key from headers
		if apiKey != "my-secret-key" {      // Replace with your actual key
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r) // Call next handler
	})
}

// LoggingMiddleware logs each request method and URL
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
		next.ServeHTTP(w, r) // Pass request to the next handler
	})
}

func experimental() {
	// from gorilla mux
	// Create a new router
	r := mux.NewRouter() // creating a router instead of multiplexer
	os.MkdirAll("uploads", os.ModePerm)

	// Apply logging middleware for requests
	r.Use(LoggingMiddleware)

	// Define routes
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// Upload route
	r.HandleFunc("/upload", UploadHandler2).Methods("POST")
	// r.HandleFunc("/progress", ProgressHandler).Methods("GET")

	// r.HandleFunc("/video/{id}", VideoHandler).Methods("GET")
	// http.Handle("/uploads/", http.StripPrefix("/uploads/", enableCORS(http.FileServer(http.Dir("uploads"))))) // to get the output.mpd

	// r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", enableCORS(http.FileServer(http.Dir("uploads/encrypted"))))) //http://localhost:8080/uploads/stream.mpd
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", enableCORS2(http.FileServer(http.Dir("uploads/encrypted")))))
	// http.Handle("/uploads/", enableCORS(http.StripPrefix("/uploads/", fs)))
	// Protected route with authentication
	// protected := r.PathPrefix("/secure").Subrouter()
	// protected.Use(AuthMiddleware)
	// protected.HandleFunc("/", ProtectedHandler).Methods("GET")

	// Start the server
	// fmt.Println("Server running on port 8080")
	// http.ListenAndServe(":3000", r)

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)

}

// Create a new ServeMux (like Express Router)
// mux := http.NewServeMux()

// // Register routes
// mux.HandleFunc("/", HomeHandler) // multiplexer
// mux.HandleFunc("/upload", UploadHandler)

// // Start the server
// fmt.Println("Server running on port 3000")
// http.ListenAndServe(":3000", mux)
