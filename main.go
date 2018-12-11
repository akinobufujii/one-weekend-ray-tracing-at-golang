package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/ftrvxmtrx/tga"
)

func main() {
	fmt.Println("start")

	width := 200
	height := 100

	outputImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Yは逆転しているので反対から書いていく
			outputImage.SetRGBA(x, height-y-1, color.RGBA{uint8(float32(x) / float32(width) * 255.99), uint8(float32(y) / float32(height) * 255.99), uint8(255 * 0.2), 255})
		}
	}

	// tgaで出力
	file, err := os.Create("result.tga")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := tga.Encode(file, outputImage); err != nil {
		panic(err)
	}

	fmt.Println("end")
}
