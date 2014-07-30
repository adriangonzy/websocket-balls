package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Ball struct {
	Position         *vector `json:"p"`
	lastGoodPosition *vector
	velocity         *vector
	Radius           float64
	mass             float64
	Color            string
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (b *Ball) String() string {
	return fmt.Sprintf("%v", &b)
}

func randFloat(min, max int) float64 {
	return rand.Float64()*float64(max-min) + float64(min)
}

func randInt(min, max int) uint {
	return uint(rand.Intn(max-min) + min)
}

func randomColor() string {
	return fmt.Sprintf("#%x", uint(rand.Float64()*float64(0xffffff)))
}

func NewRandomBall() *Ball {
	return &Ball{
		Position: &vector{randFloat(0, canvasWidth), randFloat(0, canvasHeight)},
		velocity: &vector{randFloat(0, maxVelocity), randFloat(0, maxVelocity)},
		Radius:   randFloat(0, maxRadius),
		mass:     randFloat(0, maxMass),
		Color:    randomColor(),
	}
}
