package game

import (
	"fmt"
	"time"
)

type Ball struct {
	Id       int
	Position *vector `json:"p"`
	velocity *vector
	Radius   float64
	mass     float64
	Color    string
	moved    time.Duration
}

func (b *Ball) X() float64 {
	return b.Position.X
}

func (b *Ball) Y() float64 {
	return b.Position.Y
}

func (b *Ball) String() string {
	return fmt.Sprintf("%v", &b)
}

func (b *Ball) move(delta time.Duration) {
	// convert to seconds
	acc := float64(delta/time.Millisecond) / 1000
	b.Position = b.Position.add(b.velocity.multiply(acc))
	b.moved += delta
}

func NewRandomBall(c *Config) *Ball {
	return &Ball{
		Position: &vector{randFloat(0, c.CanvasWidth/PTM), randFloat(0, c.CanvasHeight/PTM)},
		velocity: &vector{randFloat(c.MinVelocity, c.MaxVelocity), randFloat(c.MinVelocity, c.MaxVelocity)},
		Radius:   randFloat(c.MinRadius, c.MaxRadius),
		mass:     randFloat(c.MinMass, c.MaxMass),
		Color:    randomColor(),
	}
}
