package astarhexlib

import "math"

type Direction int

const (
	North Direction = iota
	NorthEast
	SouthEast
	South
	SouthWest
	NorthWest
	NumberOfDirections
)

var columns, rows int

func RotateDirection(direction Direction, amount int) Direction {
	direction = direction + Direction(amount)

	var n_dir int = int(direction) % int(NumberOfDirections)
	if n_dir < 0 {
		n_dir = int(NumberOfDirections) + n_dir
	}
	direction = Direction(n_dir)

	return direction
}

func neighbor(tile Vector2, direc Direction) Vector2 {
	if int(tile.X)%2 == 0 {
		switch direc {
		case North:
			tile.Y -= 1
		case NorthEast:
			tile.X += 1
			tile.Y -= 1
		case SouthEast:
			tile.X += 1
		case South:
			tile.Y += 1
		case SouthWest:
			tile.X -= 1
		case NorthWest:
			tile.X -= 1
			tile.Y -= 1
		}
	} else {
		switch direc {
		case North:
			tile.Y -= 1
		case NorthEast:
			tile.X += 1
		case SouthEast:
			tile.X += 1
			tile.Y += 1
		case South:
			tile.Y += 1
		case SouthWest:
			tile.X -= 1
			tile.Y += 1
		case NorthWest:
			tile.X -= 1
		}
	}

	if tile.X < 0 {
		tile.X = -1
	}
	if tile.Y < 0 {
		tile.Y = -1
	}
	if tile.X > float64(columns-1) {
		tile.X = -1
	}
	if tile.Y > float64(rows-1) {
		tile.Y = -1
	}
	return tile
}

const NODE_SIZE = 1

type ANode struct {
	Position         Vector2
	Walkable         bool
	Parent           *ANode
	DistanceToTarget float64
	Cost             float64
	Weight           float64
	F                float64
}

func NewANode(position Vector2, walkable bool) *ANode {
	return &ANode{
		Position: position,
		Walkable: walkable,
	}
}

type AStar struct {
	Grid     [][]*ANode
	GridRows int
	GridCols int
}

func NewAStar(grid [][]*ANode) *AStar {
	columns = len(grid)
	rows = len(grid[0])
	return &AStar{
		Grid:     grid,
		GridRows: len(grid[0]),
		GridCols: len(grid),
	}
}

//var Grid [][]*ANode

type Vector2 struct {
	X, Y float64
}

type PriorityQueue[T any, P float64] struct {
	items []struct {
		Element  T
		Priority P
	}
}

func NewPriorityQueue[T any, P float64]() *PriorityQueue[T, P] {
	return &PriorityQueue[T, P]{
		items: make([]struct {
			Element  T
			Priority P
		}, 0),
	}
}

func (pq *PriorityQueue[T, P]) Enqueue(element T, priority P) {
	pq.items = append(pq.items, struct {
		Element  T
		Priority P
	}{
		Element:  element,
		Priority: priority,
	})
}

func (pq *PriorityQueue[T, P]) Dequeue() T {
	minIndex := 0
	for i := 1; i < len(pq.items); i++ {
		if pq.items[i].Priority < pq.items[minIndex].Priority {
			minIndex = i
		}
	}
	item := pq.items[minIndex]
	pq.items = append(pq.items[:minIndex], pq.items[minIndex+1:]...)
	return item.Element
}

func (pq *PriorityQueue[T, P]) Count() int {
	return len(pq.items)
}

//func (a *AStar) FindPath(start, end Vector2) *Stack[*ANode] {
func (a *AStar) FindPath(start, end Vector2) []*ANode {
	startNode := NewANode(Vector2{
		X: start.X / NODE_SIZE,
		Y: start.Y / NODE_SIZE,
	}, true)
	endNode := NewANode(Vector2{
		X: end.X / NODE_SIZE,
		Y: end.Y / NODE_SIZE,
	}, true)

	//path := NewStack[*ANode]()
	var path []*ANode
	openList := NewPriorityQueue[*ANode, float64]()
	closedList := make([]*ANode, 0)

	current := startNode
	openList.Enqueue(startNode, startNode.F)

	for openList.Count() != 0 && !contains(closedList, func(n *ANode) bool {
		return n.Position == endNode.Position
	}) {
		current = openList.Dequeue()
		closedList = append(closedList, current)
		adjacencies := a.GetAdjacentNodes(current)
		for _, n := range adjacencies {
			if !contains(closedList, func(c *ANode) bool {
				return c == n
			}) && n.Walkable {
				isFound := false
				for _, oLNode := range openList.items {
					if oLNode.Element == n {
						isFound = true
						break
					}
				}
				if !isFound {
					n.Parent = current
					n.DistanceToTarget = math.Abs(n.Position.X-endNode.Position.X) + math.Abs(n.Position.Y-endNode.Position.Y)
					n.Cost = n.Weight + n.Parent.Cost
					openList.Enqueue(n, n.F)
				}
			}
		}
	}

	if !contains(closedList, func(n *ANode) bool {
		return n.Position == endNode.Position
	}) {
		return nil
	}

	temp := closedList[indexOf(closedList, func(n *ANode) bool {
		return n == current
	})]
	if temp == nil {
		return nil
	}
	for temp != startNode && temp != nil {
		//path.Push(temp)
		path = append(path, temp)
		temp = temp.Parent
	}
	return path
}

func contains[T any](slice []*T, predicate func(*T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

func indexOf[T any](slice []*T, predicate func(*T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack[T]) Count() int {
	var length int
	if s != nil {
		length = len(s.items)
		return length
	} else {
		return -1
	}
}

func (a *AStar) GetAdjacentNodes(n *ANode) []*ANode {
	var temp []*ANode
	var dir Direction

	for i := 0; i < 6; i++ {
		m := neighbor(n.Position, dir)
		if m.X > -1 && m.Y > -1 {
			temp = append(temp, a.Grid[int(m.X)][int(m.Y)])
		}
		dir = RotateDirection(dir, 1)
	}
	return temp
}
