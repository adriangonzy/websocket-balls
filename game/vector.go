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
	m := v.Magnitude()
	v.X = v.X / m
	v.Y = v.Y / m
	return v
}

func (v *vector) add(u *vector) *vector {
	return &vector{v.X + u.X, v.Y + u.Y}
}

func (v *vector) sub(u *vector) *vector {
	return &vector{v.X - u.X, v.Y - u.Y}
}

func (v *vector) multiply(l float64) *vector {
	return &vector{v.X * l, v.Y * l}
}

func (v *vector) distance(u *vector) float64 {
	d := &vector{v.X - u.X, v.Y - u.Y}
	return d.Magnitude()
}
