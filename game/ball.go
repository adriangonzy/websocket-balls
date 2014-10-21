package game

import (
	"fmt"
	"time"
)

type Ball struct {
	Id     int
	C      *vector `json:"p"`
	V      *vector `json:"_"`
	Radius float64
	Mass   float64
	Color  string
	moved  time.Duration
}

func (b *Ball) X() float64 {
	return b.C.X
}

func (b *Ball) Y() float64 {
	return b.C.Y
}

func (b *Ball) intersecting(b1 *Ball) bool {
	return b1.Radius+b.Radius > b1.C.distance(b.C)
}

func (b *Ball) String() string {
	return fmt.Sprintf("%+v", *b)
}

func (b *Ball) move(delta time.Duration) {
	// convert to seconds
	acc := float64(delta/time.Millisecond) / 1000
	b.C = b.C.add(b.V.multiply(acc))
	b.moved = b.moved + delta
}

func NewRandomBall(c *Config) *Ball {
	return &Ball{
		C:      &vector{randFloat(0, c.CanvasWidth/PTM), randFloat(0, c.CanvasHeight/PTM)},
		V:      &vector{randFloat(c.MinVelocity, c.MaxVelocity), randFloat(c.MinVelocity, c.MaxVelocity)},
		Radius: randFloat(c.MinRadius, c.MaxRadius),
		Mass:   randFloat(c.MinMass, c.MaxMass),
		Color:  randomColor(),
	}
}
