package game

import (
	"fmt"
	"sync"
	"time"
)

const (
	maxBalls     = 2
	canvasHeight = 600
	canvasWidth  = 900
	maxVelocity  = 10
	maxRadius    = 50
	maxMass      = 1000
)

var last, now time.Time
var delta time.Duration

var ballsStream chan []*ball
var balls []*ball
var wg sync.WaitGroup

func Start() chan []*ball {

	last := time.Now()

	fmt.Println("last", last)

	ballsStream := make(chan []*ball)
	balls := make([]*ball, maxBalls)

	for i := range balls {
		balls[i] = NewRandomBall()
	}

	fmt.Println(balls)

	go run()
	return ballsStream
}

func run() {
	now := time.Now()
	delta := last.Sub(now)
	last = now

	moveBalls(delta)
	wallCollisions()

	// number of pairs
	wg.Add(((len(balls) - 1) * len(balls)) / 2)

	go func() {
		wg.Wait()
		fmt.Println("===========================")
		fmt.Println(balls)
	}()

	for i, b1 := range balls {
		for _, b2 := range balls[i:] {
			go func(b1, b2 *ball) {
				if ballCollision(b1, b2) {
					collisionReaction(b1, b2)
				}
				wg.Done()
			}(b1, b2)
		}
	}
}

func moveBalls(delta time.Duration) {
	for _, b := range balls {
		b.lastGoodPosition = b.position
		b.position = b.position.add(b.velocity.multiply(uint(delta / time.Second)))
	}
}

func wallCollisions() {
	for _, b := range balls {
		if (b.position.x+uint(b.radius/2)) >= canvasWidth || (b.position.x-uint(b.radius/2)) <= 0 {
			b.position = b.lastGoodPosition
			b.position.x = -b.position.x
		}
		if (b.position.y+uint(b.radius/2)) >= canvasHeight || (b.position.y-uint(b.radius/2)) <= 0 {
			b.position = b.lastGoodPosition
			b.position.y = -b.position.y
		}
	}
}

func ballCollision(b1, b2 *ball) bool {
	return uint(b1.position.distance(b2.position)) < b1.radius+b2.radius
}

func collisionReaction(b1, b2 *ball) {
	normVector := &vector{b1.position.x - b2.position.x, b1.position.y - b2.position.y}
	normVector.Normalise()
	tangentVector := &vector{-normVector.y, normVector.x}

	b1NormalProjection := normVector.Dot(b1.velocity)
	b2NormalProjection := normVector.Dot(b2.velocity)

	b1TangentProjection := tangentVector.Dot(b1.velocity)
	b2TangentProjection := tangentVector.Dot(b2.velocity)

	// after collision
	totalMass := b1.mass + b2.mass
	b1NormalProjectionAfter := (b1NormalProjection*(b1.mass-b2.mass) + b2NormalProjection*2*b2.mass) / totalMass
	b2NormalProjectionAfter := (b2NormalProjection*(b2.mass-b1.mass) + b1NormalProjection*2*b1.mass) / totalMass

	b1.velocity = tangentVector.multiply(b1TangentProjection).add(normVector.multiply(b1NormalProjectionAfter))
	b2.velocity = tangentVector.multiply(b2TangentProjection).add(normVector.multiply(b2NormalProjectionAfter))

	b1.position = b1.lastGoodPosition
	b2.position = b2.lastGoodPosition
}
