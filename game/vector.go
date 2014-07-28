package game

import (
	"math"
)

type vector struct {
	X, Y float64
}

func (v *vector) Dot(u *vector) float64 {
	return (v.X * u.X) + (v.Y * u.Y)
}

func (v *vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v *vector) Normalise() *vector {
	v.X = v.X / v.Magnitude()
	v.Y = v.Y / v.Magnitude()
	return v
}

func (v *vector) add(u *vector) *vector {
	v.X, v.Y = v.X+u.X, v.Y+u.Y
	return v
}

func (v *vector) sub(u *vector) *vector {
	v.X, v.Y = v.X-u.X, v.Y-u.Y
	return v
}

func (v *vector) multiply(l float64) *vector {
	v.X, v.Y = v.X*l, v.Y*l
	return v
}

func (v *vector) distance(u *vector) float64 {
	d := &vector{v.X - u.X, v.Y - u.Y}
	return d.Magnitude()
}
