package utility

import (
	"math/rand"

	"github.com/ungerik/go3d/vec3"
)

// RandomInUnitSphere 単位球によるランダムな位置を生成
func RandomInUnitSphere() *vec3.T {
	sub := vec3.T{1.0, 1.0, 1.0}
	point := new(vec3.T)
	point[0] = rand.Float32()
	point[1] = rand.Float32()
	point[2] = rand.Float32()

	point.Scale(2.0)
	point.Sub(&sub)

	for vec3.Dot(point, point) >= 1.0 {
		// -1 ～ +1まで範囲での単位球内の位置を抽選する
		// 範囲外のものは再抽選
		point[0] = rand.Float32()
		point[1] = rand.Float32()
		point[2] = rand.Float32()

		point.Scale(2.0)
		point.Sub(&sub)
	}

	return point
}
