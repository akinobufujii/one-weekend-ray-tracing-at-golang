package hitable

import (
	"../material"
	"../ray"
)

// List 衝突物リスト
type List struct {
	Hitable
	HitableList []Hitable
}

// IsHit 衝突判定
func (list *List) IsHit(ray *ray.Ray, min, max float32) (isHit bool, hitResult *material.HitRecord) {
	soFar := max

	// すべての衝突判定を行う
	for _, hitable := range list.HitableList {
		if hit, hitInfo := hitable.IsHit(ray, min, soFar); hit {
			// 衝突したらmaxを小さくする
			isHit = true
			soFar = hitInfo.Distance
			hitResult = hitInfo
		}
	}

	return
}
