package camera

import (
	"../ray"
	"github.com/barnex/fmath"
	"github.com/ungerik/go3d/vec3"
)

// Camera カメラ構造体
type Camera struct {
	BottomLeft vec3.T // 左下
	Horizontal vec3.T // 水平幅
	Vertical   vec3.T // 垂直幅
	Origin     vec3.T // 中心
}

// SetParam パラメータ設定
func (camera *Camera) SetParam(pos, lockat, vup *vec3.T, fov, aspect float32) {
	halfHeight := fmath.Tan((fov * 3.14 / 180.0) * 0.5)
	halfWidth := aspect * halfHeight

	w := vec3.Sub(pos, lockat)
	w.Normalize()

	u := vec3.Cross(vup, &w)
	u.Normalize()

	v := vec3.Cross(&w, &u)
	v.Normalize()

	u.Scale(halfWidth)
	v.Scale(halfHeight)

	camera.Origin = *pos

	camera.BottomLeft = vec3.Sub(&camera.Origin, &u)
	camera.BottomLeft.Sub(&v)
	camera.BottomLeft.Sub(&w)

	camera.Horizontal = *u.Scale(2.0)
	camera.Vertical = *v.Scale(2.0)
}

// GetRay 光線獲得
func (camera *Camera) GetRay(u, v float32) *ray.Ray {
	result := new(ray.Ray)

	result.Origin = camera.Origin

	x := camera.Horizontal.Scaled(u)
	y := camera.Vertical.Scaled(v)
	result.Dir = vec3.Add(&camera.BottomLeft, &x)
	result.Dir.Add(&y)
	result.Dir.Sub(&camera.Origin)

	return result
}
