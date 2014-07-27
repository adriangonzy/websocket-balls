package game

import (
	"math"
	"testing"
)

func TestMagnitude(t *testing.T) {
	v := &vector{x: 1, y: 1}
	if v.Magnitude() != math.Sqrt(2) {
		t.Error("Expected sqrt of 2, got", v.Magnitude())
	}
}

func TestNormalise(t *testing.T) {
	v, n := &vector{x: 1, y: 1}, uint(1.0/math.Sqrt(2))
	v.Normalise()
	if v.x != n && v.y != n {
		t.Error("Expected &vector{x: 1/sqrt(2) y: 1/sqrt(2)}}, got", v)
	}
}

func TestDot(t *testing.T) {
	v := &vector{x: 1, y: 1}
	u := &vector{x: 2, y: 2}
	if v.Dot(u) != 4 {
		t.Error("Expected 4, got", v.Dot(u))
	}
}
