package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/barnex/fmath"

	"github.com/ftrvxmtrx/tga"
	"github.com/ungerik/go3d/vec3"

	"./camera"
	"./hitable"
	"./material"
	"./ray"
)

// WriteBlock 書き込みブロック
type WriteBlock struct {
	x, y int
}

// calcRayTrace 色計算（レイトレース処理）
func calcColor(ray *ray.Ray, world hitable.Hitable, depth int) vec3.T {
	// シャドウアクネ問題を解決するために
	// 極めて0に近い値を最小値として渡す
	if isHit, hitRecord := world.IsHit(ray, math.SmallestNonzeroFloat32, math.MaxFloat32); isHit {
		if isHit, attenutaion, scatteredRay := hitRecord.Material.Scatter(ray, hitRecord); depth < 50 && isHit {
			// 受け取った減数カラーを乗算しつつ、50回上限までレイトレースする
			resultColor := calcColor(scatteredRay, world, depth+1)
			return vec3.Mul(attenutaion, &resultColor)

		}

		return vec3.T{0.0, 0.0, 0.0}
	}

	t := 0.5 * (ray.Dir.Normalized()[1] + 1.0)
	return vec3.Interpolate(&vec3.T{1.0, 1.0, 1.0}, &vec3.T{0.5, 0.7, 1.0}, t)
}

// calcResultPixel 結果画素計算
func calcResultPixel(
	x, y, width, height int,
	camera *camera.Camera,
	world *hitable.List,
	outputImage *image.RGBA) {

	var calcResult vec3.T
	const samplingCount = 50
	for i := 0; i < samplingCount; i++ {
		// ジッタリングを行う
		u := (float32(x) + rand.Float32()) / float32(width)
		v := (float32(y) + rand.Float32()) / float32(height)

		// 左下からレイを飛ばして走査していく
		resultColor := calcColor(camera.GetRay(u, v), world, 0)
		calcResult.Add(&resultColor)
	}

	// ガンマ補正
	color := color.RGBA{
		uint8(fmath.Sqrt(calcResult[0]/samplingCount) * 255.99),
		uint8(fmath.Sqrt(calcResult[1]/samplingCount) * 255.99),
		uint8(fmath.Sqrt(calcResult[2]/samplingCount) * 255.99),
		255,
	}

	// Yは逆転しているので反対から書いていく
	outputImage.SetRGBA(x, height-y-1, color)
}

// calcResultPixelAsync 非同期版
func calcResultPixelAsync(
	wateGroup *sync.WaitGroup,
	block chan WriteBlock,
	width, height int,
	camera *camera.Camera,
	world *hitable.List,
	outputImage *image.RGBA) {

	defer wateGroup.Done()

	for {
		info, ok := <-block
		if !ok {
			// 処理する情報がもうない
			break
		}

		calcResultPixel(info.x, info.y, width, height, camera, world, outputImage)
	}
}

func main() {
	fmt.Println("start")

	const width = 1280
	const height = 720

	// カメラを作成して、適当なパラメータを与える
	camera := new(camera.Camera)

	camera.SetParam(
		&vec3.T{-0.5, 0.5, 1.0},
		&vec3.T{0.0, 0.0, -1.0},
		&vec3.T{0.0, 1.0, 0.0},
		90.0,
		float32(width)/float32(height))

	// レイトレース用のデータ作成
	world := new(hitable.List)

	world.HitableList = []hitable.Hitable{
		hitable.CreateSphere(&vec3.T{0.0, 0.0, -1.0}, 0.5, material.CreateLambert(&vec3.T{0.1, 0.2, 0.5})),
		hitable.CreateSphere(&vec3.T{0.0, -100.5, -1.0}, 100.0, material.CreateLambert(&vec3.T{0.8, 0.8, 0.0})),
		hitable.CreateSphere(&vec3.T{1.0, 0.0, -1.0}, 0.5, material.CreateLambert(&vec3.T{0.8, 0.6, 0.2})),
		hitable.CreateSphere(&vec3.T{-1.0, 0.0, -1.0}, 0.5, material.CreateLambert(&vec3.T{0.1, 0.2, 0.5})),
		hitable.CreateSphere(&vec3.T{-1.0, 0.0, -1.0}, -0.45, material.CreateLambert(&vec3.T{0.1, 0.2, 0.5})),
	}

	outputImage := image.NewRGBA(image.Rect(0, 0, width, height))

	rand.Seed(time.Now().UnixNano())

	const isAsync = true
	if isAsync {
		// 複数スレッド
		numCPU := runtime.NumCPU()
		ch := make(chan WriteBlock, numCPU)

		var wg sync.WaitGroup
		wg.Add(numCPU)
		runtime.GOMAXPROCS(numCPU)

		for i := 0; i < numCPU; i++ {
			go calcResultPixelAsync(&wg, ch, width, height, camera, world, outputImage)
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				ch <- WriteBlock{x, y}
			}
		}

		close(ch)
		wg.Wait()
	} else {
		// 単一スレッド
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				calcResultPixel(x, y, width, height, camera, world, outputImage)
			}
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
