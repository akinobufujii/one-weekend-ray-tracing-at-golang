package ray

import "github.com/ungerik/go3d/vec3"

// Ray 光線
type Ray struct {
	Origin vec3.T // 光線の原点
	Dir    vec3.T // 光線の方向
}

// PointAtParam 指定倍率分、光線を伸ばした位置を返す
func (ray *Ray) PointAtParam(distance float32) vec3.T {
	temp := ray.Dir.Scaled(distance)
	return vec3.Add(&ray.Origin, &temp)
}
