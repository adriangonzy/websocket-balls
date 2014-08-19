package game

import (
	"fmt"
	"sync"
	"time"
)

const (
	canvasHeight = 600 // pixels
	canvasWidth  = 900 // pixels

	PTM = 10 // pixel to meter ratio

	maxRadius   = 2  // meter
	maxVelocity = 20 // meter/s
	maxMass     = 10 // kg

	frameRate = 30                                    // frames/s
	frame     = (1000 / frameRate) * time.Millisecond // frame in ms
)

type Simulation struct {
	balls          []*Ball
	Balls          chan []*Ball
	done           chan bool
	collisions     map[int]*Collision
	collisionsChan chan *Collision
}

func makeTestBalls() []*Ball {
	balls := make([]*Ball, 2)

	balls[0] = &Ball{
		Id:       0,
		Position: &vector{10, 10},
		velocity: &vector{20, 0},
		Radius:   1,
		mass:     1,
		Color:    randomColor(),
	}

	balls[1] = &Ball{
		Id:       1,
		Position: &vector{22, 10},
		velocity: &vector{0, 0},
		Radius:   1,
		mass:     1,
		Color:    randomColor(),
	}

	return balls
}

func NewSimulation(ballCount int) *Simulation {
	fmt.Println("NEW SIMULATION")

	sim := &Simulation{}

	//balls := make([]*Ball, ballCount)
	// init random balls array
	// TODO: uniformly spread balls accross the canvas for avoiding early ball collisions
	// for i := range balls {
	// 	balls[i] = NewRandomBall()
	// 	balls[i].Id = i
	// }

	sim.balls = makeTestBalls()
	sim.Balls = make(chan []*Ball)
	sim.collisions = make(map[int]*Collision)
	sim.done = make(chan bool)

	return sim
}

func (s *Simulation) Start() {
	fmt.Println("START SIMULATION")
	ticker := time.NewTicker(frame)
	var i int
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println(i, "===============================")
				s.run(frame)
				fmt.Println(i, "===============================")
				fmt.Println()
				i++
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

			// avoid doubles
			if s.collisions[c.b1.Id] == s.collisions[c.b2.Id] {
				delete(s.collisions, c.b2.Id)
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

	fmt.Println("Moving balls for collisions")
	for _, c := range s.collisions {
		fmt.Println("collision moment", c.moment)
		c.b1.move(c.moment)
		c.b2.move(c.moment)
		collisionReaction(c.b1, c.b2)
	}

	// move balls
	convertedBalls := make([]*Ball, len(s.balls))
	for i, b := range s.balls {
		// b.wallCollision()

		b.move(delta - b.moved)
		b.moved = 0

		convertedBalls[i] = &Ball{Color: b.Color, Radius: b.Radius * PTM, Position: b.Position.multiply(PTM)}
	}

	// stream ball slice after movement computations

	s.Balls <- convertedBalls
}
