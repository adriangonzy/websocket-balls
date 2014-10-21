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

func (c *Collision) String() string {
	return fmt.Sprintf("%d", c.moment)
}

func (b *Ball) wallCollision(width, height float64) {
	r := b.Radius
	// horizontal movement collision
	switch {
	case b.Position.X+r >= width/PTM && b.velocity.X >= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = width/PTM - r
	case b.Position.X-r <= 0 && b.velocity.X <= 0:
		b.velocity.X = -b.velocity.X
		b.Position.X = r
	}

	// vertical movement collision
	switch {
	case b.Position.Y+r >= height/PTM && b.velocity.Y >= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = height/PTM - r
	case b.Position.Y-r <= 0 && b.velocity.Y <= 0:
		b.velocity.Y = -b.velocity.Y
		b.Position.Y = r
	}
}

func approaching(b1, b2 *Ball) bool {
	d1 := b1.Position.distance(b2.Position)
	b1.move(time.Millisecond)
	b2.move(time.Millisecond)
	d2 := b1.Position.distance(b2.Position)
	b1.move(-time.Millisecond)
	b2.move(-time.Millisecond)
	return d1 >= d2
}

func collisionInFrame(b1, b2 *Ball, frame time.Duration) (*Collision, bool) {

	/*
		TODO:
			- Movement before collision should be done here.
			- negative solutions have to be sorted and the most recent counted as collision
			  if balls are intersecting. test this
			- The rest of the algo will work by means of sorting and correcting movement.
	*/

	// avoid molecule like structures (glued balls)
	V1V2 := &vector{b2.velocity.X - b1.velocity.X, b2.velocity.Y - b1.velocity.Y}
	C1C2 := &vector{b2.Position.X - b1.Position.X, b2.Position.Y - b1.Position.Y}
	rTotal := b1.Radius + b2.Radius

	// discriminant computation
	a := V1V2.Dot(V1V2)
	b := 2 * C1C2.Dot(V1V2)
	c := C1C2.Dot(C1C2) - rTotal*rTotal
	D := b*b - 4*a*c

	// no possible roots
	if D < 0 {
		return nil, false
	}

	// roots computation
	t1 := (-b + math.Sqrt(D)) / (2 * a)
	t2 := (-b - math.Sqrt(D)) / (2 * a)

	// the min positive root corresponds to the collision time
	var t float64
	switch {
	case 0 < t1 && 0 < t2:
		t = math.Min(t1, t2)
	case 0 < t1 && t2 < 0:
		if b1.intersecting(b2) {
			t = t2
		}
		t = t1
	case 0 < t2 && t1 < 0:
		if b1.intersecting(b2) {
			t = t1
		}
		t = t2
	}

	// collision time in ms
	collisionTime := time.Duration(t*1000) * time.Millisecond

	if collisionTime > frame || collisionTime < -frame {
		return nil, false
	}

	return &Collision{b1, b2, collisionTime}, true
}

func (c *Collision) reaction() {
	b1, b2 := c.b1, c.b2
	//fmt.Println("COLLISION between", b1.Id, b2.Id)
	normVector := &vector{b2.Position.X - b1.Position.X, b2.Position.Y - b1.Position.Y}
	normVector.Normalise()

	// balls relative velocity projected on the normal vector
	vRelative := &vector{b2.velocity.X - b1.velocity.X, b2.velocity.Y - b1.velocity.Y}
	vRelative = normVector.multiply(vRelative.Dot(normVector))

	m1, m2 := b1.mass, b2.mass
	massMean := (m1 + m2) / 2

	// v1' = v1 + (m2/massMean) * vRelative
	b1.velocity = b1.velocity.add(vRelative.multiply(m2 / massMean))
	// v2' = v2 - (m1/massMean) * vRelative
	b2.velocity = b2.velocity.add(vRelative.multiply(-1 * m1 / massMean))
}
