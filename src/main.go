package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
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

func rgbToHsv(r, g, b uint32) (float64, float64, float64) {
	red := float64(r>>8) / 255.0
	green := float64(g>>8) / 255.0
	blue := float64(b>>8) / 255.0

	max := math.Max(red, math.Max(green, blue))
	min := math.Min(red, math.Min(green, blue))
	delta := max - min

	var hue float64

	if delta == 0 {
		hue = 0
	} else if max == red {
		hue = 60 * math.Mod((green-blue)/delta, 6)
	} else if max == green {
		hue = 60 * (((blue - red) / delta) + 2)
	} else {
		hue = 60 * (((red - green) / delta) + 4)
	}

	if hue < 0 {
		hue += 360
	}

	sat := 0.0
	if max != 0 {
		sat = delta / max
	}

	return hue, sat, max
}

func classifyColor(hue, sat, val float64) string {
	if val < 0.20 {
		return "black"
	}
	if sat < 0.20 && val > 0.80 {
		return "white"
	}
	if sat < 0.20 {
		return "gray"
	}
	if hue < 15 || hue >= 345 {
		return "red"
	} else if hue < 45 {
		return "orange"
	} else if hue < 70 {
		return "yellow"
	} else if hue < 160 {
		return "green"
	} else if hue < 250 {
		return "blue"
	} else if hue < 290 {
		return "purple"
	}
	return "pink"
}

func analyzeImages() []ImageResult {
	files, err := os.ReadDir("../img")
	if err != nil {
		fmt.Println("Kunde inte läsa img-mappen:", err)
		return nil
	}

	var results []ImageResult

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

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	results := analyzeImages()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/api/images", imagesHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
