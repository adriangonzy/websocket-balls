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
	Radius           int
	mass             int
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

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randomColor() string {
	return fmt.Sprintf("#%x", uint(rand.Float64()*float64(0xffffff)))
}

func NewRandomBall() *Ball {
	p := &vector{randInt(0, canvasWidth), randInt(0, canvasHeight)}
	return &Ball{
		Position:         p,
		lastGoodPosition: p,
		velocity:         &vector{randInt(0, maxVelocity), randInt(0, maxVelocity)},
		Radius:           randInt(1, maxRadius),
		mass:             randInt(1, maxMass),
		Color:            randomColor(),
	}
}
