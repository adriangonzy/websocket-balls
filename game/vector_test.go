package game

import (
	"math"
	"testing"
)

func TestMagnitude(t *testing.T) {
	v := &vector{X: 1, Y: 1}
	if v.Magnitude() != math.Sqrt(2) {
		t.Error("Expected sqrt of 2, got", v.Magnitude())
	}
}

func TestNormalise(t *testing.T) {
	v, n := &vector{X: 1, Y: 1}, 1.0/math.Sqrt(2)
	v.Normalise()
	if v.X != n && v.Y != n {
		t.Error("Expected &vector{X: 1/sqrt(2) Y: 1/sqrt(2)}}, got", v)
	}
}

func TestDot(t *testing.T) {
	v := &vector{X: 1, Y: 1}
	u := &vector{X: 2, Y: 2}
	if v.Dot(u) != 4 {
		t.Error("Expected 4, got", v.Dot(u))
	}
}
