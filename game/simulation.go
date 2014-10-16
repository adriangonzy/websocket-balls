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
	searchAreaFactor = PTM * 5 // should use max velocity
)

type Config struct {
	CanvasHeight float64 `json: canvasHeight` // pixels
	CanvasWidth  float64 `json: canvasWidth`  // pixels
	MaxRadius    float64 `json: maxRadius`    // meter
	MinRadius    float64 `json: minRadius`    // meter
	MaxVelocity  float64 `json: maxVelocity`  // meter/s
	MinVelocity  float64 `json: minVelocity`  // meter/s
	MaxMass      float64 `json: maxMass`      // kg
	MinMass      float64 `json: minMass`      // kg

	ballCount int   `json: numberOfBalls`
	FrameRate int   `json: framerate` // frames/s
	Frame     int64 // frame in ms
}

func (c *Config) frame() time.Duration {
	return (1000 / c.frameRate) * time.Millisecond
}

type Simulation struct {
	config     *Config
	balls      []*Ball
	Emit       chan [][]interface{}
	done       chan bool
	collisions []*Collision
}

func NewSimulation(c *Config) *Simulation {
	fmt.Println("NEW SIMULATION balls:", ballCount)

	//init random balls array
	//TODO: uniformly spread balls accross the canvas for avoiding early ball collisions
	balls := make([]*Ball, c.ballCount)
	for i := range balls {
		balls[i] = NewRandomBall()
		balls[i].Id = i
	}

	return &Simulation{balls: balls, config: c, Emit: make(chan [][]interface{}), done: make(chan bool)}
}

func (s *Simulation) Start() {
	fmt.Println("START SIMULATION")
	ticker := time.NewTicker(frame)
	var frames int
	go func() {
		for {
			select {
			case <-ticker.C:
				// TODO: block until current frame is finished
				frames = frames + 1
				fmt.Println("frame", frames)
				s.run(s.frame)
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
			b.wallCollision()

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
		canvasWidth / 2,
		canvasHeight / 2,
		canvasWidth / 2,
		canvasHeight / 2,
	}
	q := quadtree.New(box, 10)

	for _, b := range s.balls {
		q.Insert(b)
	}

	// concurrently compute pairs of balls collisions
	for _, b1 := range s.balls {
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
