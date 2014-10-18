package game

import (
	"fmt"
	"github.com/adriangonzy/websocket-balls/quadtree"
	"sort"
	"sync"
	"time"
)

const (
	PTM              = 10      // pixel to meter ratio
	searchAreaFactor = PTM * 4 // should use max velocity
)

type Config struct {
	CanvasHeight float64 `json: canvasHeight,string` // pixels
	CanvasWidth  float64 `json: canvasWidth,string`  // pixels
	MaxRadius    float64 `json: maxRadius,string`    // meter
	MinRadius    float64 `json: minRadius,string`    // meter
	MaxVelocity  float64 `json: maxVelocity,string`  // meter/s
	MinVelocity  float64 `json: minVelocity,string`  // meter/s
	MaxMass      float64 `json: maxMass,string`      // kg
	MinMass      float64 `json: minMass,string`      // kg

	BallCount int           // balls
	FrameRate int           `json: frameRate,string` // frames/s
	Frame     time.Duration // frame in ms
}

type Simulation struct {
	config     *Config
	balls      []*Ball
	Emit       chan [][]interface{}
	done       chan bool
	collisions []*Collision
}

func NewSimulation(c *Config) *Simulation {
	c.Frame = time.Duration(1000/c.FrameRate) * time.Millisecond
	fmt.Printf("NEW SIMULATION %#v", c)

	//init random balls array
	//TODO: uniformly spread balls accross the canvas for avoiding early ball collisions
	balls := make([]*Ball, c.BallCount)
	for i := range balls {
		balls[i] = NewRandomBall(c)
		balls[i].Id = i
	}

	return &Simulation{
		balls:  balls,
		config: c,
		Emit:   make(chan [][]interface{}),
		done:   make(chan bool),
	}
}

func (s *Simulation) Start() {
	fmt.Println("START SIMULATION")
	ticker := time.NewTicker(s.config.Frame)
	var frames int
	go func() {
		for {
			select {
			case <-ticker.C:
				// TODO: block until current frame is finished
				frames = frames + 1
				fmt.Println("frame", frames)
				s.run(s.config.Frame)
			case <-s.done:
				return
			}
		}
	}()
}

func (s *Simulation) Stop() {
	fmt.Println("STOP SIMULATION")
	s.done <- true
	close(s.Emit)
}

// Compute simulation balls movement during one frame
func (s *Simulation) run(delta time.Duration) {
	start := time.Now()

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
	compressedBalls := make([][]interface{}, len(s.balls))

	// finish moving balls concurrently in frame
	var wg sync.WaitGroup
	wg.Add(len(s.balls))

	for i, b := range s.balls {
		go func(b *Ball, i int) {
			// TODO: wall collision computed the same way as ball collision
			b.wallCollision(s.config.CanvasWidth, s.config.CanvasHeight)

			b.move(delta - b.moved)
			b.moved = 0

			// change to pixel unit
			p := b.Position.multiply(PTM)

			compressedBalls[i] = []interface{}{
				p.X,
				p.Y,
				b.Radius * PTM,
				b.Color,
			}

			wg.Done()
		}(b, i)
	}
	wg.Wait()

	// stream ball slice after movement computations
	s.Emit <- compressedBalls
	fmt.Println(time.Since(start))
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

	box := quadtree.Box{
		s.config.CanvasWidth / 2,
		s.config.CanvasHeight / 2,
		s.config.CanvasWidth / 2,
		s.config.CanvasHeight / 2,
	}
	q := quadtree.New(box, 10)

	for _, b := range s.balls {
		q.Insert(b)
	}

	// concurrently compute pairs of balls collisions
	for _, b1 := range s.balls {
		searchArea := s.config.MaxRadius * searchAreaFactor
		area := quadtree.Box{b1.X(), b1.Y(), searchArea, searchArea}
		// this could be optimized
		neighbors := q.SearchArea(&area)
		wg.Add(len(neighbors))
		for _, n := range neighbors {
			b2 := n.(*Ball)
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
