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
