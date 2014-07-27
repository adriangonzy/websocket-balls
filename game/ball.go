package game

import (
	"fmt"
	"math/rand"
	"time"
)

type ball struct {
	position         *vector
	lastGoodPosition *vector
	velocity         *vector
	radius           uint
	mass             uint
	color            string
}

func randInt(min, max int) uint {
	rand.Seed(time.Now().Unix())
	return uint(rand.Intn(max-min) + min)
}

func randomColor() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("#%x", uint(rand.Float32()*float32(0xffffff)))
}

func NewRandomBall() *ball {
	return &ball{
		position: &vector{randInt(0, canvasWidth), randInt(0, canvasHeight)},
		velocity: &vector{randInt(0, maxVelocity), randInt(0, maxVelocity)},
		radius:   randInt(0, maxRadius),
		mass:     randInt(0, maxMass),
		color:    randomColor(),
	}
}
