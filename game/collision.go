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
	r := float64(b.Radius)
	// horizontal movement collision
	switch {
	case b.Position.X+r >= canvasWidth && b.velocity.X >= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = canvasWidth - r
	case b.Position.X-r <= 0 && b.velocity.X <= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = r
	}

	// vertical movement collision
	switch {
	case b.Position.Y+r >= canvasHeight && b.velocity.Y >= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = canvasHeight - r
	case b.Position.Y-r <= 0 && b.velocity.Y <= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = r
	}
}

func ballCollisionInFrame(b1, b2 *Ball, delta time.Duration) (*Collision, bool) {

	C1x, C1y := b1.Position.X, b1.Position.Y
	V1x, V1y := b1.velocity.X, b1.velocity.Y
	C2x, C2y := b2.Position.X, b2.Position.Y
	V2x, V2y := b2.velocity.X, b2.velocity.Y
	r1, r2 := b1.Radius, b2.Radius

	V1V2 := V2x + V2y - V1x - V1y
	C1C2 := C2x + C2y - C1x - C1y
	rTotal := float64(r1 + r2)

	a := V1V2 * V1V2
	b := 2 * C1C2 * V1V2
	c := C1C2*C1C2 - rTotal*rTotal

	D := b*b - 4*a*c

	t1 := (-b + math.Sqrt(D)) / 2 * a
	t2 := (-b - math.Sqrt(D)) / 2 * a

	if D < 0 {
		return nil, false
	}

	var t int64
	switch {
	case 0 < t1 && 0 < t2:
		t = int64(math.Min(t1, t2))
	case 0 < t1 && t2 < 0:
		t = int64(t1)
	case 0 < t2 && t1 < 0:
		t = int64(t2)
	}

	frame := int64(delta / time.Millisecond)

	fmt.Println("=========================")
	fmt.Println("D", D)
	fmt.Println("t1", t1)
	fmt.Println("t2", t2)
	fmt.Println("t min", t)
	fmt.Println("delta", frame)

	if t > frame || t <= 0 {
		fmt.Println("NO COLLISION IN FRAME")
		return nil, false
	}

	fmt.Println("COLLISION IN FRAME")
	return &Collision{b1, b2, time.Duration(t) * time.Millisecond}, true
}

func collisionReaction(b1, b2 *Ball) {
	normVector := &vector{b1.Position.X - b2.Position.X, b1.Position.Y - b2.Position.Y}
	normVector.Normalise()
	tangentVector := &vector{-normVector.Y, normVector.X}

	b1NormalProjection := normVector.Dot(b1.velocity)
	b2NormalProjection := normVector.Dot(b2.velocity)

	b1TangentProjection := tangentVector.Dot(b1.velocity)
	b2TangentProjection := tangentVector.Dot(b2.velocity)

	// after collision
	m1, m2 := float64(b1.mass), float64(b2.mass)
	totalMass := m1 + m2
	b1NormalProjectionAfter := ((m1-m2)/totalMass)*b1NormalProjection + (2*m2/totalMass)*b2NormalProjection
	b2NormalProjectionAfter := ((m2-m1)/totalMass)*b2NormalProjection + (2*m1/totalMass)*b1NormalProjection

	b1.velocity = tangentVector.multiply(b1TangentProjection).add(normVector.multiply(b1NormalProjectionAfter))
	b2.velocity = tangentVector.multiply(b2TangentProjection).add(normVector.multiply(b2NormalProjectionAfter))
}
