package game

import (
	"fmt"
	"testing"
)

func TestRandomColor(t *testing.T) {
	fmt.Println(randomColor())
}

func TestNewRandomBall(t *testing.T) {
	fmt.Println("%v", NewRandomBall())
}
