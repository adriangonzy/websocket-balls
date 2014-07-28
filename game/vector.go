package game

import (
	"math"
)

type vector struct {
	x, y uint
}

func (v *vector) Dot(u *vector) uint {
	return (v.x * u.x) + (v.y * u.y)
}

func (v *vector) Magnitude() float64 {
	return math.Sqrt(float64(v.Dot(v)))
}

func (v *vector) Normalise() *vector {
	v.x = uint(float64(v.x) / v.Magnitude())
	v.y = uint(float64(v.y) / v.Magnitude())
	return v
}

func (v *vector) add(u *vector) *vector {
	v.x, v.y = v.x+u.x, v.y+u.y
	return v
}

func (v *vector) sub(u *vector) *vector {
	v.x, v.y = v.x-u.x, v.y-u.y
	return v
}

func (v *vector) multiply(l uint) *vector {
	v.x, v.y = v.x*l, v.y*l
	return v
}

func (v *vector) distance(u *vector) float64 {
	d := &vector{v.x - u.x, v.y - u.y}
	return d.Magnitude()
}
