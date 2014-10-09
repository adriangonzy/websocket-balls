package game

import (
	"testing"
)

func TestStartSimulation(t *testing.T) {
	NewSimulation(2).Start()
}

func makeTestBalls() []*Ball {
	balls := make([]*Ball, 2)

	balls[0] = &Ball{
		Id:       1,
		Position: &vector{10, 10},
		velocity: &vector{20, 0},
		Radius:   randFloat(1, 4),
		mass:     4,
		Color:    randomColor(),
	}

	balls[1] = &Ball{
		Id:       2,
		Position: &vector{22, 10},
		velocity: &vector{-13, 0},
		Radius:   randFloat(1, 7),
		mass:     1,
		Color:    randomColor(),
	}

	return balls
}
