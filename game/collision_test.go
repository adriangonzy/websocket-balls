package game

import (
	"fmt"
	"testing"
)

func TestBallCollisionInFrame(t *testing.T) {
	b1 := &Ball{
		Id:       1,
		Position: &vector{10, 10},
		velocity: &vector{5, 0},
		Radius:   1,
		mass:     1,
		Color:    randomColor(),
	}

	b2 := &Ball{
		Id:       2,
		Position: &vector{22, 10},
		velocity: &vector{0, 0},
		Radius:   1,
		mass:     1,
		Color:    randomColor(),
	}

	if c, ok := ballCollisionInFrame(b1, b2, frame); ok {
		fmt.Println("collision moment", c.moment)
		b1.move(c.moment)
		b2.move(c.moment)
	}

}
