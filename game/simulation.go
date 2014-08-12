package game

import (
	"fmt"
	"math"
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

type Collision struct {
	b1, b2 *Ball
	moment time.Duration
}

type Simulation struct {
	balls []*Ball
	Balls chan []*Ball
	done  chan bool
	collisions []*Collision
	collisionsChan chan *Collision
}

func NewSimulation(ballCount int) *Simulation {
	sim := &Simulation{}
	balls := make([]*Ball, ballCount)

	// init random balls array
	// TODO: uniformly spread balls accross the canvas for avoiding early early ball collisions
	for i := range balls {
		balls[i] = NewRandomBall()
		balls[i].Id = i
	}

	sim.balls = balls
	sim.Balls = make(chan []*Ball)
	sim.collisions = make([]*Collision, 1)
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
	var wg sync.WaitGroup

	// number of ball pairs
	wg.Add(((len(s.balls) - 1) * len(s.balls)) / 2)

	// init collision reception channel
	s.collisionsChan = make(chan *Collision)
	go func() {
		for c := range s.collisionsChan {
			if len(s.collisions) == 0 {
				s.collisions = append(s.collisions, collision)
			} else {
				for i, c := range s.collisions {
					// if collision happens before another collision
					// insert new collision
					if collision.moment < c.moment {
						s.collisions = append(s.collisions[:i], append([]*Collision{collision}, s.collisions[i:]...)...)
						break		
					}
				}	
			}
		}
	}()

	// concurrently compute pairs of balls collisions
	for i, b1 := range s.balls {
		for j, b2 := range s.balls[i+1:] {
			go func(b1, b2 *Ball) {
				if collision, ok := ballCollisionInFrame(b1, b2); ok {
					s.collisionsChan <- collision
				}
				wg.Done()
			}(b1, b2)
		}
	}

	go func() {
		wg.Wait()
		close(s.collisionsChan)
		for _, b := range s.balls {

			// use collisions slice to compute frame intermediate movements
					
			b.move(delta)
			b.wallCollision()
		}
	}()

	// stream ball slice after movement computations
	s.Balls <- s.balls
}

func (b *Ball) move(delta time.Duration) {
	b.Position = b.Position.add(b.velocity.multiply(float64(delta/time.Millisecond) / 10))
}

func (b *Ball) wallCollision() {
	r := float64(b.Radius)
	// horizontal movement collision
	switch {
	case b.Position.X+r >= canvasWidth && b.velocity.X >= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = canvasWidth - r
	case b.Position.X-r <= 0 && b.velocity.X <= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = r
	}

	// vertical movement collision
	switch {
	case b.Position.Y+r >= canvasHeight && b.velocity.Y >= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = canvasHeight - r
	case b.Position.Y-r <= 0 && b.velocity.Y <= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = r
	}
}

func ballCollisionInFrame(b1, b2 *Ball) bool, time.Duration {
	return int(b1.Position.distance(b2.Position)) < b1.Radius+b2.Radius
}

func collisionReaction(b1, b2 *Ball) {
	normVector := &vector{b1.Position.X - b2.Position.X, b1.Position.Y - b2.Position.Y}
	normVector.Normalise()
	tangentVector := &vector{-normVector.Y, normVector.X}

	d := b1.Position.distance(b2.Position)
	rT := float64(b1.Radius + b2.Radius)
	c := (rT - d) / math.Sqrt(2.0)

	fmt.Println("before")
	fmt.Println("b1", b1.Position)
	fmt.Println("b2", b2.Position)
	fmt.Println("T radius", b1.Radius+b2.Radius)
	fmt.Println("distance", b1.Position.distance(b2.Position))

	fmt.Println("correction coeff", c)
	c1 := normVector.multiply(c)
	c2 := normVector.multiply(-1 * c)
	b1.Position = b1.Position.add(c1)
	b2.Position = b2.Position.add(c2)

	fmt.Println("after correction")
	fmt.Println("c1", c1)
	fmt.Println("c2", c2)
	fmt.Println("b1", b1.Position)
	fmt.Println("b2", b2.Position)
	fmt.Println("T radius", b1.Radius+b2.Radius)
	fmt.Println("distance", b1.Position.distance(b2.Position))
	fmt.Println("=================")

	b1NormalProjection := normVector.Dot(b1.velocity)
	b2NormalProjection := normVector.Dot(b2.velocity)

	b1TangentProjection := tangentVector.Dot(b1.velocity)
	b2TangentProjection := tangentVector.Dot(b2.velocity)

	// after collision
	m1, m2 := float64(b1.mass), float64(b2.mass)
	totalMass := m1 + m2
	b1NormalProjectionAfter := ((m1-m2)/totalMass)*b1NormalProjection + (2*m2/totalMass)*b2NormalProjection
	b2NormalProjectionAfter := ((m2-m1)/totalMass)*b2NormalProjection + (2*m1/totalMass)*b1NormalProjection

	b1.velocity = tangentVector.multiply(b1TangentProjection).add(normVector.multiply(b1NormalProjectionAfter))
	b2.velocity = tangentVector.multiply(b2TangentProjection).add(normVector.multiply(b2NormalProjectionAfter))
}
