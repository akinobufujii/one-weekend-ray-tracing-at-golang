package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/ftrvxmtrx/tga"
	"github.com/ungerik/go3d/vec3"

	"./camera"
	"./ray"
)

// calcRayTrace 色計算（レイトレース処理）
func calcRayTrace(ray *ray.Ray) vec3.T {
	a := vec3.T{1.0, 1.0, 1.0}
	b := vec3.T{0.5, 0.7, 1.0}
	t := 0.5 * (ray.Dir.Normalized()[1] + 1.0)

	return vec3.Interpolate(&a, &b, t)
}

func main() {
	fmt.Println("start")

	width := 200
	height := 100

	// カメラを作成して、適当なパラメータを与える
	camera := new(camera.Camera)

	pos := vec3.T{-0.5, 0.5, 1.0}
	lookat := vec3.T{0.0, 0.0, -1.0}
	vup := vec3.T{0.0, 1.0, 0.0}
	camera.SetParam(&pos, &lookat, &vup, 90.0, float32(width)/float32(height))

	outputImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			ray := camera.GetRay(float32(x)/float32(width), float32(y)/float32(height))

			// レイトレース計算
			calcResult := calcRayTrace(ray)

			// ガンマ補正
			color := color.RGBA{}
			color.R = uint8(calcResult[0] * 255.99)
			color.G = uint8(calcResult[1] * 255.99)
			color.B = uint8(calcResult[2] * 255.99)
			color.A = 255

			// Yは逆転しているので反対から書いていく
			outputImage.SetRGBA(x, height-y-1, color)
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
