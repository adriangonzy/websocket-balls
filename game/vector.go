package game

import (
	"math"
)

type vector struct {
	X, Y int
}

func (v *vector) Dot(u *vector) int {
	return (v.X * u.X) + (v.Y * u.Y)
}

func (v *vector) Magnitude() float64 {
	return math.Sqrt(float64(v.Dot(v)))
}

func (v *vector) Normalise() *vector {
	v.X = int(float64(v.X) / v.Magnitude())
	v.Y = int(float64(v.Y) / v.Magnitude())
	return v
}

func (v *vector) add(u *vector) *vector {
	return &vector{v.X + u.X, v.Y + u.Y}
}

func (v *vector) sub(u *vector) *vector {
	return &vector{v.X - u.X, v.Y - u.Y}
}

func (v *vector) multiply(l float64) *vector {
	return &vector{int(float64(v.X) * l), int(float64(v.Y) * l)}
}

func (v *vector) distance(u *vector) float64 {
	d := &vector{v.X - u.X, v.Y - u.Y}
	return d.Magnitude()
}
