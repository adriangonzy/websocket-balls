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
	balls          []*Ball
	Balls          chan []*Ball
	done           chan bool
	collisions     map[int]*Collision
	collisionsChan chan *Collision
}

func NewSimulation(ballCount int) *Simulation {
	fmt.Println("NEW SIMULATION")

	sim := &Simulation{}
	balls := make([]*Ball, ballCount)

	// init random balls array
	// TODO: uniformly spread balls accross the canvas for avoiding early ball collisions
	for i := range balls {
		balls[i] = NewRandomBall()
		balls[i].Id = i
		fmt.Println(i)
		fmt.Println(balls[i])
	}

	sim.balls = balls
	sim.Balls = make(chan []*Ball)
	sim.collisions = make(map[int]*Collision)
	sim.done = make(chan bool)

	return sim
}

func (s *Simulation) Start() {
	fmt.Println("START SIMULATION")
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
	fmt.Println("STOP SIMULATION")
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

			if s.collisions[c.b1.Id] == nil {
				s.collisions[c.b1.Id] = c
			}

			if s.collisions[c.b2.Id] == nil {
				s.collisions[c.b2.Id] = c
			}

			b1M := s.collisions[c.b1.Id].moment
			b2M := s.collisions[c.b2.Id].moment
			cM := c.moment

			if b1M < b2M {
				delete(s.collisions, c.b2.Id)
				if cM < b1M {
					s.collisions[c.b1.Id] = c
				}
			} else if b2M < b1M {
				delete(s.collisions, c.b1.Id)
				if cM < b2M {
					s.collisions[c.b2.Id] = c
				}
			}
		}
	}()

	// concurrently compute pairs of balls collisions
	for i, b1 := range s.balls {
		for _, b2 := range s.balls[i+1:] {
			go func(b1, b2 *Ball) {
				if collision, ok := ballCollisionInFrame(b1, b2, delta); ok {
					s.collisionsChan <- collision
				}
				wg.Done()
			}(b1, b2)
		}
	}

	wg.Wait()
	close(s.collisionsChan)

	// move balls
	for _, b := range s.balls {
		b.wallCollision()

		for _, c := range s.collisions {
			c.b1.move(c.moment)
			c.b2.move(c.moment)
			collisionReaction(c.b1, c.b2)
		}

		b.move(delta - b.moved)
		b.moved = 0
	}

	// stream ball slice after movement computations
	s.Balls <- s.balls
}
