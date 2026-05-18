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
	files, err := os.ReadDir("../public/img")
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

		file, err := os.Open("../public/img/" + name)
		if err != nil {
			continue
		}

		img, _, err := image.Decode(file)
		file.Close()

		if err != nil {
			continue
		}

		colorCount := make(map[string]int)
		hsvTracker := make(map[string][]float64)

		bounds := img.Bounds()

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

		for color, count := range colorCount {
			if count > maxPixels {
				maxPixels = count
				category = color
			}
		}

		hue := hsvTracker[category][0]
		sat := hsvTracker[category][1]
		val := hsvTracker[category][2]

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
	http.Handle("/", http.FileServer(http.Dir("../public")))
	http.HandleFunc("/api/images", imagesHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
