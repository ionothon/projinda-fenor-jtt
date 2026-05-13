package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

func main() {
	file, err := os.Open("../img/test2.jpg")
	if err != nil {
		fmt.Println("Hittade inte bilden: ", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Kunde inte läsa bilden: ", err)
		return
	}

	redBucket := 0
	greenBucket := 0
	blueBucket := 0

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			R := r >> 8
			G := g >> 8
			B := b >> 8

			if R > G && R > B {
				redBucket++
			} else if G > R && G > B {
				greenBucket++
			} else if B > R && B > G {
				blueBucket++
			}
		}
	}
	fmt.Printf("Resultat: Röd - %d, Grön - %d, Blå - %d", redBucket, greenBucket, blueBucket)

}
