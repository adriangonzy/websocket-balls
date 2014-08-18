package game

import (
	"fmt"
	"time"
)

type Ball struct {
	Id               int
	Position         *vector `json:"p"`
	lastGoodPosition *vector
	velocity         *vector
	Radius           int
	mass             int
	Color            string
	moved            time.Duration
}

func (b *Ball) String() string {
	return fmt.Sprintf("%v", &b)
}

func (b *Ball) move(delta time.Duration) {
	b.Position = b.Position.add(b.velocity.multiply(float64(delta/time.Millisecond) / 10))
	b.moved += delta
}

func NewRandomBall() *Ball {
	p := &vector{randFloat(0, canvasWidth), randFloat(0, canvasHeight)}
	return &Ball{
		Position:         p,
		lastGoodPosition: p,
		velocity:         &vector{randFloat(0, maxVelocity), randFloat(0, maxVelocity)},
		Radius:           randInt(1, maxRadius),
		mass:             randInt(1, maxMass),
		Color:            randomColor(),
	}
}
