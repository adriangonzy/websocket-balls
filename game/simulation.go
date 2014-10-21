package game

import (
	"fmt"
	"github.com/adriangonzy/websocket-balls/quadtree"
	"sort"
	"sync"
	"time"
)

const (
	PTM = 10 // pixel to meter ratio
)

type Config struct {
	CanvasHeight     float64 `json: canvasHeight` // pixels
	CanvasWidth      float64 `json: canvasWidth`  // pixels
	MaxRadius        float64 `json: maxRadius`    // meter
	MinRadius        float64 `json: minRadius`    // meter
	MaxVelocity      float64 `json: maxVelocity`  // meter/s
	MinVelocity      float64 `json: minVelocity`  // meter/s
	MaxMass          float64 `json: maxMass`      // kg
	MinMass          float64 `json: minMass`      // kg
	FrameRate        int     `json: frameRate`    // frames/s
	SearchAreaFactor int     `json: searchAreaFactor`
	BallCount        int

	Frame time.Duration // frame in ms
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
	fmt.Printf("NEW SIMULATION %#v\n", c)

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
				fmt.Println("===================")
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

func print(msg string) {
	fmt.Println(msg)
}

// Compute simulation balls movement during one frame
func (s *Simulation) run(delta time.Duration) {
	start := time.Now()
	s.computeCollisions(delta)
	fmt.Println("collisions", len(s.collisions), "time", time.Since(start))

	fmt.Printf("%#v\n", s.balls)
	fmt.Println("compute")
	if len(s.collisions) > 0 {
		s.sortCollisions()
		fmt.Printf("%#v\n", s.balls)
		fmt.Println("sort")
		s.moveAfterCollisions()
		fmt.Printf("%#v\n", s.balls)
		fmt.Println("move after")
	}
	s.finishMoving(delta)
	fmt.Printf("%#v\n", s.balls)
	fmt.Println("finish")

	// stream ball slice after movement computations
	s.Emit <- s.compressBalls()
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
		searchArea := s.config.MaxRadius * float64(s.config.SearchAreaFactor) * PTM
		area := quadtree.Box{b1.C.X, b1.C.Y, searchArea, searchArea}
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

func (s *Simulation) moveAfterCollisions() {
	collided := make(map[int]bool)
	// compute ball movement given collision slice
	for _, c := range s.collisions {
		if collided[c.B1.Id] || collided[c.B2.Id] {
			continue
		}
		collided[c.B1.Id], collided[c.B2.Id] = true, true
		// move balls to collision time
		c.B1.move(c.moment)
		c.B2.move(c.moment)
		c.reaction()
	}
}

func (s *Simulation) finishMoving(delta time.Duration) {
	// finish moving balls concurrently in frame
	var wg sync.WaitGroup
	wg.Add(len(s.balls))
	for i, b := range s.balls {
		go func(b *Ball, i int) {
			// TODO: wall collision computed the same way as ball collision
			b.wallCollision(s.config.CanvasWidth, s.config.CanvasHeight)
			b.move(delta - b.moved)
			b.moved = 0
			wg.Done()
		}(b, i)
	}
	wg.Wait()
}

func (s *Simulation) compressBalls() [][]interface{} {
	// change to pixel unit and compress sent data
	compressedBalls := make([][]interface{}, len(s.balls))
	fmt.Printf("%#v\n", s.balls)
	for i, b := range s.balls {
		p := b.C.multiply(PTM)
		compressedBalls[i] = []interface{}{
			p.X,
			p.Y,
			b.Radius * PTM,
			b.Color,
		}
	}
	return compressedBalls
}

type ByTime []*Collision

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].moment < a[j].moment }

func (s *Simulation) sortCollisions() {

	sort.Sort(ByTime(s.collisions))
}
