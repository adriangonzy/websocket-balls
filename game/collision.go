package game

import (
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

	// TODO

	return nil, false
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
