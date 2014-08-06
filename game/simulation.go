package game

import (
	"fmt"
	"sync"
	"time"
)

const (
	canvasHeight = 600
	canvasWidth  = 900
	maxVelocity  = 7
	maxRadius    = 50
	maxMass      = 5
	frameRate    = 30
	frameTimer   = 1000 / frameRate
)

type Simulation struct {
	balls []*Ball
	Balls chan []*Ball
	done  chan bool
}

func NewSimulation(ballCount int) *Simulation {
	sim := &Simulation{}
	balls := make([]*Ball, ballCount)

	// init random balls array
	// TODO: uniformly spread balls accross the canvas for avoiding early early ball collisions
	for i := range balls {
		balls[i] = NewRandomBall()
	}

	sim.balls = balls
	sim.Balls = make(chan []*Ball)
	sim.done = make(chan bool)

	return sim
}

func (s *Simulation) Start() {
	ticker := time.NewTicker(time.Millisecond * frameTimer)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.run(frameTimer * time.Millisecond)
			case <-s.done:
				return
			}
		}
	}()
}

func (s *Simulation) Stop() {
	s.done <- true
	close(s.Balls)
}

// Compute simulation balls movement during one frame
func (s *Simulation) run(delta time.Duration) {

	s.moveBalls(delta)
	s.wallCollisions()

	var wg sync.WaitGroup

	// number of ball pairs
	wg.Add(((len(s.balls) - 1) * len(s.balls)) / 2)
	go func() {
		wg.Wait()
	}()

	// concurrently compute pairs of balls collisions
	for i, b1 := range s.balls {
		for j, b2 := range s.balls[i+1:] {
			go func(b1, b2 *Ball) {
				if ballCollision(b1, b2) {
					fmt.Printf("Collision %v : %v \n", i, j)
					collisionReaction(b1, b2)
				}
				wg.Done()
			}(b1, b2)
		}
	}

	// stream ball slice after movement computations
	s.Balls <- s.balls
}

func (s *Simulation) moveBalls(delta time.Duration) {
	for _, b := range s.balls {
		b.lastGoodPosition = b.Position
		dMovement := b.velocity.multiply(float64(delta/time.Millisecond) / 10)
		b.Position = b.Position.add(dMovement)
	}
}

func (s *Simulation) wallCollisions() {
	for _, b := range s.balls {
		if b.Position.X+b.Radius >= canvasWidth || b.Position.X-b.Radius <= 0 {
			fmt.Println("vertical collision ball", b)
			b.Position = b.lastGoodPosition
			b.velocity.X = -b.velocity.X
		}
		if b.Position.Y+b.Radius >= canvasHeight || b.Position.Y-b.Radius <= 0 {
			fmt.Println("horizontal collision ball", b)
			b.Position = b.lastGoodPosition
			b.velocity.Y = -b.velocity.Y
		}
	}
}

func ballCollision(b1, b2 *Ball) bool {
	return int(b1.Position.distance(b2.Position)) < b1.Radius+b2.Radius
}

func collisionReaction(b1, b2 *Ball) {
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

	fmt.Println("TOTAL MASS FOR COLLISION", totalMass)

	b1.velocity = tangentVector.multiply(float64(b1TangentProjection)).add(normVector.multiply(float64(b1NormalProjectionAfter)))
	b2.velocity = tangentVector.multiply(float64(b2TangentProjection)).add(normVector.multiply(float64(b2NormalProjectionAfter)))

	b1.Position = b1.lastGoodPosition
	b2.Position = b2.lastGoodPosition
}
