package entity

import (
	"github.com/gemrs/gem/gem/game/interface/entity"
	"github.com/gemrs/gem/gem/game/position"
)

/* mob directions */
const (
	DirectionNorth     int = 1
	DirectionNorthEast int = 2
	DirectionEast      int = 4
	DirectionSouthEast int = 7
	DirectionSouth     int = 6
	DirectionSouthWest int = 5
	DirectionWest      int = 3
	DirectionNorthWest int = 0
	DirectionNone      int = -1
)

var directionMap = [3][3]int{
	{DirectionNorthEast, DirectionEast, DirectionSouthEast},
	{DirectionNorth, DirectionNone, DirectionSouth},
	{DirectionNorthWest, DirectionWest, DirectionSouthWest},
}

// A SimpleWaypointQueue trusts the points generated by the client and does
// simple interpolation to determine the next position.
// In future, we might want to create another implementation which performs
// server-side path finding
type SimpleWaypointQueue struct {
	points        []*position.Absolute
	lastDirection int
	direction     int
}

func NewSimpleWaypointQueue() *SimpleWaypointQueue {
	return &SimpleWaypointQueue{
		points:        make([]*position.Absolute, 0),
		lastDirection: DirectionNone,
		direction:     DirectionNone,
	}
}

// Empty determines if there are any points queued
func (q *SimpleWaypointQueue) Empty() bool {
	return len(q.points) == 0
}

// Clear clears the waypoint queue
func (q *SimpleWaypointQueue) Clear() {
	q.points = []*position.Absolute{}
}

// Push appends a point to the waypoint queue
func (q *SimpleWaypointQueue) Push(point *position.Absolute) {
	q.points = append(q.points, point)
}

// Tick advances the waypoint queue, and returns the next position of the mob
func (q *SimpleWaypointQueue) Tick(mob entity.Movable) {
	if len(q.points) == 0 {
		// Nothing to do
		return
	}

	current := mob.Position()
	nextWaypoint := q.points[0]
	if current.Compare(nextWaypoint) {
		// We've reached a waypoint, dequeue it and continue
		q.points = q.points[1:]
		q.Tick(mob)
		return
	}

	next := current.NextInterpolatedPoint(nextWaypoint)

	dx, dy, _ := current.Delta(next)

	q.lastDirection = q.direction
	q.direction = directionMap[dx+1][dy+1]

	mob.SetNextStep(next)
}

// WalkDirection returns the mob's current and (in the case of running) last walking direction
func (q *SimpleWaypointQueue) WalkDirection() (current, last int) {
	return q.direction, q.lastDirection
}
