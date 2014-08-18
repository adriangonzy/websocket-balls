package game

import (
	_ "fmt"
	"testing"
	"time"
)

func TestBallCollisionInFrame(t *testing.T) {
	b1, b2 := NewRandomBall(), NewRandomBall()
	ballCollisionInFrame(b1, b2, frameTimer*time.Millisecond)
}
