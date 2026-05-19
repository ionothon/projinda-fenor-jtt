package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

// ImageResult stores analyzed information about image
type ImageResult struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Hue      float64 `json:"hue"`
	Sat      float64 `json:"sat"`
	Val      float64 `json:"val"`
}

// analyzeImages reads images from the img folder,
// analyzes their colors and applies filters.
func analyzeImages(colorFilter, toneFilter string) []ImageResult {
	// Read all files inside the image folder.
	files, err := os.ReadDir("../public/img")
	if err != nil {
		fmt.Println("Could not read img folder:", err)
		return nil
	}

	results := []ImageResult{}

	for _, entry := range files {
		// Skip directories.
		if entry.IsDir() {
			continue
		}

		// Skip non-image files
		name := entry.Name()
		if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".jpeg") && !strings.HasSuffix(name, ".png") {
			continue
		}

		file, err := os.Open("../public/img/" + name)
		if err != nil {
			fmt.Printf("Something went wrong decoding the file %s: %v\n", name, err)
			continue
		}

		img, _, err := image.Decode(file)
		file.Close()

		if err != nil {
			continue
		}
		// colorCount counts how many pixels belong to each color category
		colorCount := make(map[string]int)
		// hsvTracker stores the HSV values for each color category
		hsvTracker := make(map[string][]float64)

		// Gets image dimensions
		bounds := img.Bounds()

		// Iterates through each pixel, converting RGB to HSV and categorizing it while increasing the count
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()

				h, s, v := rgbToHsv(r, g, b)
				cat := classifyColor(h, s, v)
				colorCount[cat]++
				hsvTracker[cat] = []float64{h, s, v}
			}
		}
		category := "Unknown"
		maxPixels := 0
		// Finds the dominant color by comparing pixel counts for each category
		for color, count := range colorCount {
			if count > maxPixels {
				maxPixels = count
				category = color
			}
		}
		// Gets HSV values for the dominant color
		hue := hsvTracker[category][0]
		sat := hsvTracker[category][1]
		val := hsvTracker[category][2]

		// Determine cool or warm tone
		tone := getTone(hue)

		// Skip image if it does not match selected color filter
		if colorFilter != "" && !strings.Contains(colorFilter, category) {
			continue
		}
		// Skip image if it does not match selected tone filter
		if toneFilter != "" && tone != toneFilter {
			continue
		}

		results = append(results, ImageResult{
			Name:     name,
			Category: category,
			Hue:      hue,
			Sat:      sat,
			Val:      val,
		})
	}

	return results
}

func getTone(hue float64) string {
	if hue < 70 || hue >= 290 {
		return "warm"
	}
	return "cool"
}

// imagesHandler handles requests to the image API, reading query filters in URL and converting response to JSON
func imagesHandler(response http.ResponseWriter, request *http.Request) {
	colorFilter := request.URL.Query().Get("color")
	toneFilter := request.URL.Query().Get("tone")
	results := analyzeImages(colorFilter, toneFilter)

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(results)
}

// uploadHandler handles image uploads from the frontend
func uploadHandler(response http.ResponseWriter, request *http.Request) {

	// Only allow POST requests for uploads
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded form data.
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	// Get uploaded image files from the form.
	files := request.MultipartForm.File["uploadedImages"]

	// Loop through uploaded files.
	for _, fileHeader := range files {

		// Open uploaded file.
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create destination file in image folder.
		destinationFile, err := os.Create("../public/img/" + fileHeader.Filename)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		defer destinationFile.Close()

		// Copy uploaded file data into destination file.
		_, err = io.Copy(destinationFile, file)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response.WriteHeader(http.StatusOK)

	// Return success response after upload.
	response.Write([]byte("Upload successful!"))
}

// Register frontend, API endpoints and start local server.
func main() {

	// Serve frontend files from the public folder.
	http.Handle("/", http.FileServer(http.Dir("../public")))

	// Register API endpoint for image analysis.
	http.HandleFunc("/api/images", imagesHandler)

	// Register API endpoint for image uploads.
	http.HandleFunc("/api/upload", uploadHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
