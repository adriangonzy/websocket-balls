package game

import (
	"fmt"
	"math"
	"time"
)

type Collision struct {
	B1, B2 *Ball
	moment time.Duration
}

func (c *Collision) String() string {
	return fmt.Sprintf("%+v", *c)
}

func (b *Ball) wallCollision(width, height float64) {
	r := b.Radius
	// horizontal movement collision
	switch {
	case b.C.X+r >= width/PTM && b.V.X >= 0:
		b.V.X = -b.V.X
		b.C.X = width/PTM - r
	case b.C.X-r <= 0 && b.V.X <= 0:
		b.V.X = -b.V.X
		b.C.X = r
	}

	// vertical movement collision
	switch {
	case b.C.Y+r >= height/PTM && b.V.Y >= 0:
		b.V.Y = -b.V.Y
		b.C.Y = height/PTM - r
	case b.C.Y-r <= 0 && b.V.Y <= 0:
		b.V.Y = -b.V.Y
		b.C.Y = r
	}
}

func collisionInFrame(b1, b2 *Ball, frame time.Duration) (*Collision, bool) {

	if b1.Id == b2.Id {
		return nil, false
	}

	// avoid molecule like structures (glued balls)
	V1V2 := &vector{b2.V.X - b1.V.X, b2.V.Y - b1.V.Y}
	C1C2 := &vector{b2.C.X - b1.C.X, b2.C.Y - b1.C.Y}
	rTotal := b1.Radius + b2.Radius

	TFrame := float64(frame) / float64(time.Millisecond*1000)

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
		if b1.intersecting(b2) && t2 > -TFrame {
			t = t2
		}
		t = t1
	case 0 < t2 && t1 < 0:
		if b1.intersecting(b2) && t1 > -TFrame {
			t = t1
		}
		t = t2
	}

	// collision time in ms
	collisionTime := time.Duration(t*1000) * time.Millisecond

	if collisionTime > frame || collisionTime < -frame || t == 0 {
		return nil, false
	}

	//* debug
	if b1.intersecting(b2) {
		fmt.Println(b1.Id, b2.Id, "intersecting", "t1", t1, "t2", t2)
	} else {
		fmt.Println(b1.Id, b2.Id, "t1", t1, "t2", t2)
	}

	fmt.Println(b1.Id, "collides", b2.Id, "at", collisionTime)
	//*/ debug

	return &Collision{b1, b2, collisionTime}, true
}

func (c *Collision) reaction() {
	b1, b2 := c.B1, c.B2

	//fmt.Println("COLLISION between", b1.Id, b2.Id)
	normVector := &vector{b2.C.X - b1.C.X, b2.C.Y - b1.C.Y}
	normVector.Normalise()

	// balls relative velocity projected on the normal vector
	vRelative := &vector{b2.V.X - b1.V.X, b2.V.Y - b1.V.Y}
	vRelative = normVector.multiply(vRelative.Dot(normVector))

	m1, m2 := b1.Mass, b2.Mass
	massMean := (m1 + m2) / 2

	// v1' = v1 + (m2/massMean) * vRelative
	b1.V = b1.V.add(vRelative.multiply(m2 / massMean))
	// v2' = v2 - (m1/massMean) * vRelative
	b2.V = b2.V.add(vRelative.multiply(-1 * m1 / massMean))
}
