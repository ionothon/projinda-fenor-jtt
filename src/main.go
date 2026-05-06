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

	pixel := img.At(50, 120)
	r, g, b, _ := pixel.RGBA()

	fmt.Printf("Röd: %d, Grön: %d, Blå: %d\n", r>>8, g>>8, b>>8)

}
