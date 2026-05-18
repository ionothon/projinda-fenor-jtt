package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strings"
)

type ImageResult struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Hue      float64 `json:"hue"`
	Sat      float64 `json:"sat"`
	Val      float64 `json:"val"`
}

func analyzeImages(colorFilter, toneFilter string) []ImageResult {
	files, err := os.ReadDir("../img")
	if err != nil {
		fmt.Println("Kunde inte läsa img-mappen:", err)
		return nil
	}

	results := []ImageResult{}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".jpeg") && !strings.HasSuffix(name, ".png") {
			continue
		}

		file, err := os.Open("../img/" + name)
		if err != nil {
			continue
		}

		img, _, err := image.Decode(file)
		file.Close()

		if err != nil {
			continue
		}

		bounds := img.Bounds()
		centerX := bounds.Max.X / 2
		centerY := bounds.Max.Y / 2

		r, g, b, _ := img.At(centerX, centerY).RGBA()
		hue, sat, val := rgbToHsv(r, g, b)
		category := classifyColor(hue, sat, val)
		tone := getTone(hue)

		if colorFilter != "" && category != colorFilter {
			continue
		}
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

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	colorFilter := r.URL.Query().Get("color")
	toneFilter := r.URL.Query().Get("tone")
	results := analyzeImages(colorFilter, toneFilter)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/api/images", imagesHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
