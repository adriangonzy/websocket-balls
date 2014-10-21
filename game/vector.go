package game

import (
	"fmt"
	"math"
)

type vector struct {
	X, Y float64
}

func (v *vector) String() string {
	return fmt.Sprintf("%+v", *v)
}

func (v *vector) Dot(u *vector) float64 {
	return (v.X * u.X) + (v.Y * u.Y)
}

func (v *vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v *vector) Normalise() *vector {
	m := math.Sqrt(v.Dot(v))
	if m == 0 {
		panic("Impossible to divide by zero")
	}
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
	t := &vector{v.X * l, v.Y * l}
	return t
}

func (v *vector) distance(u *vector) float64 {
	d := &vector{v.X - u.X, v.Y - u.Y}
	return d.Magnitude()
}
