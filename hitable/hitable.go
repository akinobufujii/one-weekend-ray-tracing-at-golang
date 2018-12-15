package hitable

import (
	"../material"
	"../ray"
)

// Hitable 衝突物インターフェイス
type Hitable interface {
	IsHit(ray *ray.Ray, min, max float32) (isHit bool, hitResult *material.HitRecord)
}
