package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("../img/test.jpg")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	fmt.Println("Opened")
}
