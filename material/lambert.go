package material

import (
	"math/rand"

	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/ray"
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/utility"
	"github.com/ungerik/go3d/vec3"
)

// Lambert 拡散反射マテリアル
type Lambert struct {
	Material
	Albedo vec3.T // アルベドカラー
}

// CreateLambert 拡散反射マテリアル作成
func CreateLambert(albedo *vec3.T) *Lambert {
	lambert := &Lambert{}
	lambert.Albedo = *albedo
	return lambert
}

// Scatter 散乱結果を返す
func (lambert *Lambert) Scatter(randomDevice *rand.Rand, in *ray.Ray, hitRecord *HitRecord) (isHit bool, attenuation *vec3.T, scattered *ray.Ray) {
	// 衝突したところから更にランダムな位置に光線を飛ばして
	// 各種結果を返す
	target := hitRecord.Point
	target.Add(&hitRecord.Normal)
	target.Add(utility.RandomInUnitSphere(randomDevice))

	scattered = new(ray.Ray)
	scattered.Origin = hitRecord.Point
	scattered.Dir = vec3.Sub(&target, &hitRecord.Point)
	attenuation = &lambert.Albedo

	isHit = true
	return
}
