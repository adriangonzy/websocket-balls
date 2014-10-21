// Package quadtree implements methods for a quadtree spatial partitioning data
// structure.
//
// Code is based on the Wikipedia article
// http://en.wikipedia.org/wiki/Quadtree.
package quadtree

import (
	_ "fmt"
)

// XY is a simple coordinate structure for points and vectors
type Point interface {
	X() float64
	Y() float64
}

// AABB represents an Axis-Aligned bounding box structure with center and half
// dimension
type Box struct {
	CenterX float64
	CenterY float64
	HalfX   float64
	HalfY   float64
}

// NewAABB creates a new axis-aligned bounding box and returns its address
func NewBox(CenterX, CenterY, HalfX, HalfY float64) *Box {
	return &Box{CenterX, CenterY, HalfX, HalfY}
}

// ContainsPoint returns true when the AABB contains the point given
func (b *Box) ContainsPoint(p Point) bool {
	if p.X() < b.CenterX-b.HalfX {
		return false
	}
	if p.X() > b.CenterX+b.HalfX {
		return false
	}
	if p.Y() < b.CenterY-b.HalfY {
		return false
	}
	if p.Y() > b.CenterY+b.HalfY {
		return false
	}

	return true
}

// IntersectsAABB returns true when the AABB intersects another AABB
func (b *Box) IntersectsBox(other *Box) bool {
	if other.CenterX+other.HalfX < b.CenterX-b.HalfX {
		return false
	}
	if other.CenterY+other.HalfY < b.CenterY-b.HalfY {
		return false
	}
	if other.CenterX-other.HalfX > b.CenterX+b.HalfX {
		return false
	}
	if other.CenterY-other.HalfY > b.CenterY+b.HalfY {
		return false
	}
	return true
}

// QuadTree represents the quadtree data structure.
type QuadTree struct {
	boundary     Box
	points       []Point
	nodeCapacity int
	northWest    *QuadTree
	northEast    *QuadTree
	southWest    *QuadTree
	southEast    *QuadTree
}

// New creates a new quadtree node that is bounded by boundary and contains
// nodeCapacity points.
// nodeCapacity is the maximum number of points allowed in a quadtree node
func New(boundary Box, nodeCapacity int) *QuadTree {
	points := make([]Point, 0, nodeCapacity)
	qt := &QuadTree{
		boundary:     boundary,
		points:       points,
		nodeCapacity: nodeCapacity,
	}
	return qt
}

// Insert adds a point to the quadtree. It returns true if it was successful
// and false otherwise.
func (qt *QuadTree) Insert(p Point) bool {
	// Ignore objects which do not belong in this quad tree.
	if !qt.boundary.ContainsPoint(p) {
		return false
	}

	// If there is space in this quad tree, add the object here.
	if len(qt.points) < cap(qt.points) {
		qt.points = append(qt.points, p)
		return true
	}

	// Otherwise, we need to subdivide then add the point to whichever node
	// will accept it.
	if qt.northWest == nil {
		qt.subDivide()
	}

	if qt.northWest.Insert(p) {
		return true
	}
	if qt.northEast.Insert(p) {
		return true
	}
	if qt.southWest.Insert(p) {
		return true
	}
	if qt.southEast.Insert(p) {
		return true
	}

	// Otherwise, the point cannot be inserted for some unknown reason.
	// (which should never happen)
	return false
}

func (qt *QuadTree) subDivide() {
	// Check if this is a leaf node.
	if qt.northWest != nil {
		return
	}

	box := Box{
		qt.boundary.CenterX - qt.boundary.HalfX/2,
		qt.boundary.CenterY + qt.boundary.HalfY/2,
		qt.boundary.HalfX / 2,
		qt.boundary.HalfY / 2,
	}
	qt.northWest = New(box, qt.nodeCapacity)

	box = Box{
		qt.boundary.CenterX + qt.boundary.HalfX/2,
		qt.boundary.CenterY + qt.boundary.HalfY/2,
		qt.boundary.HalfX / 2,
		qt.boundary.HalfY / 2,
	}
	qt.northEast = New(box, qt.nodeCapacity)

	box = Box{
		qt.boundary.CenterX - qt.boundary.HalfX/2,
		qt.boundary.CenterY - qt.boundary.HalfY/2,
		qt.boundary.HalfX / 2,
		qt.boundary.HalfY / 2,
	}
	qt.southWest = New(box, qt.nodeCapacity)

	box = Box{
		qt.boundary.CenterX + qt.boundary.HalfX/2,
		qt.boundary.CenterY - qt.boundary.HalfY/2,
		qt.boundary.HalfX / 2,
		qt.boundary.HalfY / 2,
	}
	qt.southEast = New(box, qt.nodeCapacity)

	for _, v := range qt.points {
		if qt.northWest.Insert(v) {
			continue
		}
		if qt.northEast.Insert(v) {
			continue
		}
		if qt.southWest.Insert(v) {
			continue
		}
		if qt.southEast.Insert(v) {
			continue
		}
	}
	qt.points = nil
}

func (q *QuadTree) isLeaf() {
	return qt.northWest == nil
}

func (q *QuadTree) leafs() []*QuadTree {
	leafs := []*QuadTree{}

	if !q.isLeaf() {
		leafs = append(leafs, qt.northWest.leafs())
		leafs = append(leafs, qt.northEast.leafs())
		leafs = append(leafs, qt.southWest.leafs())
		leafs = append(leafs, qt.southEast.leafs())
	} else {
		leafs = append(leafs, q)
	}
	return leafs
}

func (q *QuadTree) leafPoints() [][]Point {
	leafs := q.leafs()
	points := [][]Point{}
	for _, l := range leafs {
		points = append(points, l.points)
	}
	return points
}

func (qt *QuadTree) SearchArea(a *Box) []Point {
	results := make([]Point, 0, qt.nodeCapacity)

	if !qt.boundary.IntersectsBox(a) {
		return results
	}

	for _, v := range qt.points {
		if a.ContainsPoint(v) {
			results = append(results, v)
		}
	}

	if qt.northWest == nil {
		return results
	}

	results = append(results, qt.northWest.SearchArea(a)...)
	results = append(results, qt.northEast.SearchArea(a)...)
	results = append(results, qt.southWest.SearchArea(a)...)
	results = append(results, qt.southEast.SearchArea(a)...)
	return results
}
