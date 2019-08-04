package hitable

import (
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/material"
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/ray"
	"github.com/barnex/fmath"
	"github.com/ungerik/go3d/vec3"
)

// Sphere 球
type Sphere struct {
	Hitable
	Center   vec3.T
	Radius   float32
	Material material.Material
}

// CreateSphere 球作成
func CreateSphere(center *vec3.T, radius float32, material material.Material) (sphere *Sphere) {
	sphere = new(Sphere)
	sphere.Center = *center
	sphere.Radius = radius
	sphere.Material = material
	return
}

// IsHit 衝突判定
func (sphere *Sphere) IsHit(ray *ray.Ray, min, max float32) (isHit bool, hitResult *material.HitRecord) {
	centerToRay := vec3.Sub(&ray.Origin, &sphere.Center)

	// 半径の2乗と(光線 - 球の中心)の2乗は等しい
	// (光線 - 球の中心)の2乗 = (光線 - 球の中心)・(光線 - 球の中心)と表すことができる
	// そのため、「(光線 - 球の中心)・(光線 - 球の中心) = 半径の2乗」と表せる
	a := vec3.Dot(&ray.Dir, &ray.Dir)
	b := vec3.Dot(&centerToRay, &ray.Dir)
	c := vec3.Dot(&centerToRay, &centerToRay) - (sphere.Radius * sphere.Radius)

	// 光線と球があたっているかの解の公式に当てはめる
	// →「x = (-b +- sqrt(b * b - a * c)) / a」
	// 平方根内の計算が実数(0より上)ならあたっている可能性がある
	discreminat := (b * b) - (a * c)
	if discreminat <= 0.0 {
		// 0以下ならあたってない
		isHit = false
		return
	}

	// 判定して、結果を書き込む無名関数
	judgeAndWriteResult := func(distance float32) bool {
		if min < distance && distance < max {
			// あたったので情報を書き込む
			hitResult = new(material.HitRecord)
			hitResult.Distance = distance
			hitResult.Point = ray.PointAtParam(distance)
			hitResult.Normal = vec3.Sub(&hitResult.Point, &sphere.Center)
			hitResult.Normal[0] /= sphere.Radius
			hitResult.Normal[1] /= sphere.Radius
			hitResult.Normal[2] /= sphere.Radius
			hitResult.Material = sphere.Material
			return true
		}

		// あたってない
		return false
	}

	// 球の中心から光線が出ている方向にあたっているか判定
	if judgeAndWriteResult((-b - fmath.Sqrt(discreminat)) / a) {
		isHit = true
		return
	}

	// 反対側判定
	if judgeAndWriteResult((-b + fmath.Sqrt(discreminat)) / a) {
		isHit = true
		return
	}

	// ここまできたらあたってない
	isHit = false
	return
}
