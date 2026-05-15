package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"math"
	"os"
)

func rgbToHsv(r, g, b uint32) (float64, float64, float64) {
	red := float64(r>>8) / 255.0
	green := float64(g>>8) / 255.0
	blue := float64(b>>8) / 255.0

	max := math.Max(red, math.Max(green, blue))
	min := math.Max(red, math.Min(green, blue))
	delta := max - min

	var hue, sat, val float64

	if delta == 0 {
		hue = 0
	} else if max == red {
		hue = 60 * math.Mod((green-blue)/delta, 6)
	} else if max == green {
		hue = 60 * ((blue - red) + 2)
	} else if max == blue {
		hue = 60 * ((red - green) + 4)
	}

	if hue < 0 {
		hue += 360
	}

	if max == 0 {
		sat = 0
	} else {
		sat = delta / max
	}

	val = max

	return hue, sat, val
}

func main() {

	files, err := os.ReadDir("../img")
	if err != nil {
		fmt.Println("Kunde inte läsa img-mappen:", err)
		return
	}
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		file, err := os.Open("../img/" + entry.Name())
		if err != nil {
			fmt.Println("Hittade inte bilden:", entry.Name(), err)
			continue
		}
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			fmt.Println("Kunde inte läsa bilden:", entry.Name(), err)
			continue
		}
		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				hue, sat, val := rgbToHsv(r, g, b)
				if x == bounds.Max.X/2 && y == bounds.Max.Y/2 {
					fmt.Printf("%s -> Hue: %.2f, Sat: %.2f, Val %.2f\n", entry.Name(), hue, sat, val)
				}
			}
		}
	}
}
