package game

import (
	"fmt"
	"math"
	"time"
)

type Collision struct {
	b1, b2 *Ball
	moment time.Duration
}

func (b *Ball) wallCollision() {
	r := b.Radius
	// horizontal movement collision
	switch {
	case b.Position.X+r >= canvasWidth/PTM && b.velocity.X >= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = canvasWidth/PTM - r
	case b.Position.X-r <= 0 && b.velocity.X <= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = r
	}

	// vertical movement collision
	switch {
	case b.Position.Y+r >= canvasHeight/PTM && b.velocity.Y >= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = canvasHeight/PTM - r
	case b.Position.Y-r <= 0 && b.velocity.Y <= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = r
	}
}

func ballCollisionInFrame(b1, b2 *Ball, frame time.Duration) (*Collision, bool) {

	V1V2 := b2.velocity.add(b1.velocity.multiply(-1))
	C1C2 := b2.Position.add(b1.Position.multiply(-1))
	rTotal := b1.Radius + b2.Radius

	a := V1V2.Dot(V1V2)
	b := 2 * C1C2.Dot(V1V2)
	c := C1C2.Dot(C1C2) - rTotal*rTotal

	D := b*b - 4*a*c

	if D < 0 {
		return nil, false
	}

	t1 := (-b + math.Sqrt(D)) / (2 * a)
	t2 := (-b - math.Sqrt(D)) / (2 * a)

	var t float64
	switch {
	case 0 < t1 && 0 < t2:
		t = math.Min(t1, t2)
	case 0 < t1 && t2 < 0:
		t = t1
	case 0 < t2 && t1 < 0:
		t = t2
	}

	moment := time.Duration(t*1000) * time.Millisecond

	if moment > frame || t <= 0 {
		return nil, false
	}

	return &Collision{b1, b2, moment}, true
}

func collisionReaction(b1, b2 *Ball) {
	fmt.Println("COLLISION between", b1.Id, b2.Id)
	normVector := &vector{b2.Position.X - b1.Position.X, b2.Position.Y - b1.Position.Y}
	normVector.Normalise()

	vRelative := &vector{b2.velocity.X - b1.velocity.X, b2.velocity.Y - b1.velocity.Y}
	vRelative = normVector.multiply(vRelative.Dot(normVector))

	m1, m2 := b1.mass, b2.mass
	massMean := (m1 + m2) / 2

	b1.velocity = b1.velocity.add(vRelative.multiply(m2 / massMean))
	b2.velocity = b2.velocity.add(vRelative.multiply(-1 * m1 / massMean))
}
