package game

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	canvasHeight = 600 // pixels
	canvasWidth  = 900 // pixels

	PTM = 10 // pixel to meter ratio

	maxRadius   = 0.6 // meter
	minRadius   = 0.6 // meter
	maxVelocity = 25  // meter/s
	minVelocity = -25 // meter/s
	maxMass     = 10  // kg
	minMass     = 1   // kg

	frameRate = 30                                    // frames/s
	frame     = (1000 / frameRate) * time.Millisecond // frame in ms
)

type Simulation struct {
	balls      []*Ball
	Balls      chan []*Ball
	done       chan bool
	collisions []*Collision
}

func NewSimulation(ballCount int) *Simulation {
	fmt.Println("NEW SIMULATION balls:", ballCount)

	//init random balls array
	//TODO: uniformly spread balls accross the canvas for avoiding early ball collisions
	balls := make([]*Ball, ballCount)
	for i := range balls {
		balls[i] = NewRandomBall()
		balls[i].Id = i
	}

	return &Simulation{balls: balls, Balls: make(chan []*Ball), done: make(chan bool)}
}

func (s *Simulation) Start() {
	fmt.Println("START SIMULATION")
	ticker := time.NewTicker(frame)
	go func() {
		for {
			select {
			case <-ticker.C:
				// TODO: block until current frame is finished
				s.run(frame)
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

	s.computeCollisions(delta)
	s.sortCollisions()

	collided := make(map[*Ball]bool)

	// compute ball movement given collision slice
	for _, c := range s.collisions {
		if collided[c.b1] || collided[c.b2] {
			continue
		}
		collided[c.b1], collided[c.b2] = true, true
		c.reaction()
	}

	// change to pixel unit
	convertedBalls := make([]*Ball, len(s.balls))

	// finish moving balls concurrently in frame
	var wg sync.WaitGroup
	wg.Add(len(s.balls))

	for i, b := range s.balls {
		go func(b *Ball, i int) {
			// TODO: wall collision computed the same way as ball collision
			b.wallCollision()

			b.move(delta - b.moved)
			b.moved = 0

			// change to pixel unit
			convertedBalls[i] = &Ball{
				Color:    b.Color,
				Radius:   b.Radius * PTM,
				Position: b.Position.multiply(PTM),
			}
			wg.Done()
		}(b, i)
	}
	wg.Wait()

	// stream ball slice after movement computations
	s.Balls <- convertedBalls
}

func (s *Simulation) computeCollisions(delta time.Duration) {
	// clear past collisions
	s.collisions = []*Collision{}

	// init collision reception channel
	cols := make(chan *Collision)
	go func() {
		for c := range cols {
			s.collisions = append(s.collisions, c)
		}
	}()

	// number of ball pairs
	var wg sync.WaitGroup
	wg.Add(((len(s.balls) - 1) * len(s.balls)) / 2)

	// concurrently compute pairs of balls collisions
	// TODO: use quad tree
	for i, b1 := range s.balls {
		for _, b2 := range s.balls[i+1:] {
			go func(b1, b2 *Ball) {
				if c, ok := collisionInFrame(b1, b2, delta); ok {
					cols <- c
				}
				wg.Done()
			}(b1, b2)
		}
	}

	wg.Wait()
	close(cols)
}

type ByTime []*Collision

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].moment < a[j].moment }

func (s *Simulation) sortCollisions() {
	sort.Sort(ByTime(s.collisions))
}
