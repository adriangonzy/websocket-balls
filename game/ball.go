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

func (b *Ball) String() string {
	return fmt.Sprintf("%v", &b)
}

func (b *Ball) move(delta time.Duration) {
	// convert to seconds
	acc := float64(delta/time.Millisecond) / 1000
	b.Position = b.Position.add(b.velocity.multiply(acc))
	b.moved += delta
}

func NewRandomBall() *Ball {
	return &Ball{
		Position: &vector{randFloat(0, canvasWidth/PTM), randFloat(0, canvasHeight/PTM)},
		velocity: &vector{randFloat(minVelocity, maxVelocity), randFloat(minVelocity, maxVelocity)},
		Radius:   randFloat(minRadius, maxRadius),
		mass:     randFloat(minMass, maxMass),
		Color:    randomColor(),
	}
}
