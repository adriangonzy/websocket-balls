package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	maxBalls     = 3
	canvasHeight = 600
	canvasWidth  = 900
	maxVelocity  = 10
	maxRadius    = 50
	maxMass      = 1000
	frameRate    = 60
	frameTimer   = 1000 / frameRate
)

var last, now time.Time
var delta time.Duration

var ballsStream chan []*ball
var balls []*ball
var wg sync.WaitGroup

func Start() chan []*ball {

	last = time.Now()

	ballsStream = make(chan []*ball)
	balls = make([]*ball, maxBalls)

	for i := range balls {
		balls[i] = NewRandomBall()
	}

	ticker := time.NewTicker(time.Millisecond * frameTimer)
	go func() {
		for {
			<-ticker.C
			run(frameTimer * time.Millisecond)
		}
	}()

	return ballsStream
}

func formatBalls() []interface{} {
	bs := make([]interface{}, len(balls)*4)
	for i, j := 0, 0; i < len(balls) && j < len(balls)*4; i++ {
		bs[j] = balls[i].Position.X
		bs[j+1] = balls[i].Position.Y
		bs[j+2] = balls[i].Radius
		bs[j+3] = balls[i].Color
		j += 4
	}
	return bs
}

func run(delta time.Duration) {

	moveBalls(delta)
	wallCollisions()

	j, e := json.Marshal(formatBalls())
	if e != nil {
		fmt.Errorf("%v", e)
	}
	fmt.Printf("%s\n", j)

	wg.Add(((len(balls) - 1) * len(balls)) / 2)

	go func() {
		wg.Wait()
	}()

	for i, b1 := range balls {
		for j, b2 := range balls[i+1:] {
			go func(b1, b2 *ball) {
				if ballCollision(b1, b2) {
					fmt.Printf("Collision %v : %v \n", i, j)
					collisionReaction(b1, b2)
				}
				wg.Done()
			}(b1, b2)
		}
	}
	//ballsStream <- balls
}

func moveBalls(delta time.Duration) {
	for _, b := range balls {
		b.lastGoodPosition = b.Position
		dMovement := b.velocity.multiply(float64(delta/time.Millisecond) / 1000)
		b.Position = b.Position.add(dMovement)
	}
}

func wallCollisions() {
	for _, b := range balls {
		if b.Position.X+b.Radius >= canvasWidth || b.Position.X-b.Radius <= 0 {
			b.Position = b.lastGoodPosition
			b.Position.X = -b.Position.X
		}
		if b.Position.Y+b.Radius >= canvasHeight || b.Position.Y-b.Radius <= 0 {
			b.Position = b.lastGoodPosition
			b.Position.Y = -b.Position.Y
		}
	}
}

func ballCollision(b1, b2 *ball) bool {
	return b1.Position.distance(b2.Position) < b1.Radius+b2.Radius
}

func collisionReaction(b1, b2 *ball) {
	normVector := &vector{b1.Position.X - b2.Position.X, b1.Position.Y - b2.Position.Y}
	normVector.Normalise()
	tangentVector := &vector{-normVector.Y, normVector.X}

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

	b1.Position = b1.lastGoodPosition
	b2.Position = b2.lastGoodPosition
}
