package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

func main() {

	files, err := os.ReadDir("../img")
	if err != nil {
		fmt.Println("Could not read img folder:", err)
		return
	}
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		file, err := os.Open("../img/" + entry.Name())
		if err != nil {
			fmt.Println("Could not find image:", entry.Name(), err)
			continue
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			fmt.Println("Could not read image:", entry.Name(), err)
			continue
		}

		colorCount := make(map[string]int)

		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()

				hue, sat, val := rgbToHsv(r, g, b)
				category := hsvSorter(hue, sat, val)
				colorCount[category]++
			}
		}
		dominantColor := "Unknown"
		maxPixels := 0

		for color, count := range colorCount {
			if count > maxPixels {
				maxPixels = count
				dominantColor = color
			}
		}
		fmt.Printf("Bild: %s Dominant färg: %s (%d pixlar)\n", entry.Name(), dominantColor, maxPixels)
	}
}
