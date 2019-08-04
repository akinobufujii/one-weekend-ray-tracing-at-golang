package hitable

import (
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/material"
	"github.com/akinobufujii/one-weekend-ray-tracing-at-golang/ray"
)

// Hitable 衝突物インターフェイス
type Hitable interface {
	IsHit(ray *ray.Ray, min, max float32) (isHit bool, hitResult *material.HitRecord)
}
