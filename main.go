package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/barnex/fmath"

	"github.com/ungerik/go3d/vec3"

	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/camera"
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/hitable"
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/material"
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/ray"
)

// WriteLine 書き込みライン
type WriteLine struct {
	y, width int
}

// ReadLine 読み込み込みライン
type ReadLine struct {
	y      int
	colors []color.RGBA
}

// calcRayTrace 色計算（レイトレース処理）
func calcColor(randomDevice *rand.Rand, ray *ray.Ray, world hitable.Hitable, depth int) vec3.T {
	// シャドウアクネ問題を解決するために
	// 極めて0に近い値を最小値として渡す
	if isHit, hitRecord := world.IsHit(ray, math.SmallestNonzeroFloat32, math.MaxFloat32); isHit {
		if isHit, attenutaion, scatteredRay := hitRecord.Material.Scatter(randomDevice, ray, hitRecord); depth < 50 && isHit {
			// 受け取った減数カラーを乗算しつつ、50回上限までレイトレースする
			resultColor := calcColor(randomDevice, scatteredRay, world, depth+1)
			return vec3.Mul(attenutaion, &resultColor)
		}

		return vec3.T{0.0, 0.0, 0.0}
	}

	t := 0.5 * (ray.Dir.Normalized()[1] + 1.0)
	return vec3.Interpolate(&vec3.T{1.0, 1.0, 1.0}, &vec3.T{0.5, 0.7, 1.0}, t)
}

// calcResultPixel 結果画素計算
func calcResultPixel(
	randomDevice *rand.Rand,
	x, y, imageWidth, imageHeight int,
	camera *camera.Camera,
	world *hitable.List) color.RGBA {

	var calcResult vec3.T
	const samplingCount = 50
	for i := 0; i < samplingCount; i++ {
		// ジッタリングを行う
		u := (float32(x) + randomDevice.Float32()) / float32(imageWidth)
		v := (float32(y) + randomDevice.Float32()) / float32(imageHeight)

		// 左下からレイを飛ばして走査していく
		resultColor := calcColor(randomDevice, camera.GetRay(u, v), world, 0)
		calcResult.Add(&resultColor)
	}

	// ガンマ補正
	return color.RGBA{
		uint8(fmath.Sqrt(calcResult[0]/samplingCount) * 255.99),
		uint8(fmath.Sqrt(calcResult[1]/samplingCount) * 255.99),
		uint8(fmath.Sqrt(calcResult[2]/samplingCount) * 255.99),
		255,
	}
}

func main() {
	fmt.Println("start")

	const imageWidth = 1920
	const imageHeight = 1080

	// カメラを作成して、適当なパラメータを与える
	camera := new(camera.Camera)

	camera.SetParam(
		&vec3.T{-0.5, 0.5, 1.0},
		&vec3.T{0.0, 0.0, -1.0},
		&vec3.T{0.0, 1.0, 0.0},
		90.0,
		float32(imageWidth)/float32(imageHeight))

	// レイトレース用のデータ作成
	world := new(hitable.List)

	world.HitableList = []hitable.Hitable{
		hitable.CreateSphere(&vec3.T{0.0, 0.0, -1.0}, 0.5, material.CreateLambert(&vec3.T{0.1, 0.2, 0.5})),
		hitable.CreateSphere(&vec3.T{0.0, -100.5, -1.0}, 100.0, material.CreateLambert(&vec3.T{0.8, 0.8, 0.0})),
		hitable.CreateSphere(&vec3.T{1.0, 0.0, -1.0}, 0.5, material.CreateLambert(&vec3.T{0.8, 0.6, 0.2})),
		hitable.CreateSphere(&vec3.T{-1.0, 0.0, -1.0}, 0.5, material.CreateLambert(&vec3.T{0.1, 0.2, 0.5})),
		hitable.CreateSphere(&vec3.T{-1.0, 0.0, -1.0}, -0.45, material.CreateLambert(&vec3.T{0.1, 0.2, 0.5})),
	}

	outputImage := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	const isAsync = true
	if isAsync {
		// 複数スレッド
		numCPU := runtime.NumCPU()

		writeLineStream := make(chan WriteLine, numCPU)
		readLineStream := make(chan ReadLine, numCPU)
		doneSendReadLine := make(chan struct{}, numCPU)
		doneAll := make(chan struct{})

		// 情報送信goroutine
		go func() {
			// 横1列を高さ分渡して並列化
			for y := 0; y < imageHeight; y++ {
				writeLineStream <- WriteLine{y, imageWidth}
			}
			close(writeLineStream)
		}()

		for i := 0; i < numCPU; i++ {
			// ピクセル計算goroutine
			go func(randomDevice *rand.Rand) {
				for info := range writeLineStream {
					color := make([]color.RGBA, info.width)
					for x := 0; x < info.width; x++ {
						color[x] = calcResultPixel(randomDevice, x, info.y, imageWidth, imageHeight, camera, world)
					}
					readLineStream <- ReadLine{info.y, color}
				}
				doneSendReadLine <- struct{}{}
			}(rand.New(rand.NewSource(time.Now().Unix())))
		}

		// ピクセル合成goroutine
		go func() {
			for readLine := range readLineStream {
				for x, color := range readLine.colors {
					// Yは逆転しているので反対から書いていく
					outputImage.SetRGBA(x, imageHeight-readLine.y-1, color)
				}
			}
			doneAll <- struct{}{}
		}()

		for i := 0; i < numCPU; i++ {
			<-doneSendReadLine
		}
		close(readLineStream)
		<-doneAll
	} else {
		// 単一スレッド
		randomDevice := rand.New(rand.NewSource(time.Now().Unix()))
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				color := calcResultPixel(randomDevice, x, y, imageWidth, imageHeight, camera, world)

				// Yは逆転しているので反対から書いていく
				outputImage.SetRGBA(x, imageHeight-y-1, color)
			}
		}
	}

	// pngで出力
	file, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, outputImage); err != nil {
		panic(err)
	}

	fmt.Println("end")
}
