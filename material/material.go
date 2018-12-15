package material

import (
	"../ray"
	"github.com/ungerik/go3d/vec3"
)

// Material マテリアルインターフェイス
type Material interface {
	// Scatter 散乱結果を返す
	Scatter(in *ray.Ray, hitRecord *HitRecord) (isHit bool, attenuation *vec3.T, scattered *ray.Ray)
}

// HitRecord 衝突結果
type HitRecord struct {
	Distance float32
	Point    vec3.T
	Normal   vec3.T
	Material Material
}
