package main

import (
	"container/heap"

	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Chiton"
}

type point struct {
	x int
	y int
}

type node struct {
	p        point
	distance int
	index    int
}

type priorityQueue []*node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueue) update(node *node, distance int) {
	node.distance = distance
	heap.Fix(pq, node.index)
}

// frontier wraps a priority queue and also tracks the
// memory location of each node by its x/y coordinates
// for easy lookup.
type frontier struct {
	queue priorityQueue
	nodes map[point]*node
	size  int
}

func newFrontier(points map[point]int, size int) frontier {
	infinity := 9*len(points) + 1
	f := frontier{
		queue: make(priorityQueue, len(points)),
		nodes: map[point]*node{},
		size:  size,
	}

	i := 0
	for point := range points {
		n := node{
			p:        point,
			distance: infinity,
			index:    i,
		}
		f.queue[i] = &n
		f.nodes[point] = &n
		i++
	}
	heap.Init(&f.queue)

	return f
}

func (f *frontier) update(p point, distance int) {
	n := f.nodes[p]
	f.queue.update(n, distance)
}

func (f *frontier) pop() *node {
	return heap.Pop(&f.queue).(*node)
}

func (f *frontier) isEmpty() bool {
	return f.queue.Len() <= 0
}

func (f *frontier) get(p point) *node {
	return f.nodes[p]
}

func (f *frontier) neighbors(p point) []*node {
	var nearby []*node
	if p.x > 0 {
		nearby = append(nearby, f.nodes[point{p.x - 1, p.y}])
	}
	if p.y > 0 {
		nearby = append(nearby, f.nodes[point{p.x, p.y - 1}])
	}
	if p.x < f.size-1 {
		nearby = append(nearby, f.nodes[point{p.x + 1, p.y}])
	}
	if p.y < f.size-1 {
		nearby = append(nearby, f.nodes[point{p.x, p.y + 1}])
	}
	return nearby
}

func dijkstra(points map[point]int, size int) int {
	begin := point{0, 0}
	end := point{size - 1, size - 1}

	frontier := newFrontier(points, size)
	frontier.update(begin, 0)

	for !frontier.isEmpty() {
		current := frontier.pop()

		for _, neighbor := range frontier.neighbors(current.p) {
			proposed := current.distance + points[neighbor.p]
			if proposed < neighbor.distance {
				frontier.update(neighbor.p, proposed)
			}
		}

		if current.p == end {
			break
		}
	}

	return frontier.get(end).distance
}

func (s *solver) Solve1(input []string) (int, error) {
	points := map[point]int{}
	size := len(input)
	for y, line := range input {
		ints := aoc.IntCharacters(line)
		for x, risk := range ints {
			points[point{x, y}] = risk
		}
	}

	return dijkstra(points, size), nil
}

func explode(points map[point]int, size int) map[point]int {
	exploded := map[point]int{}
	for p, risk := range points {
		for xFactor := 0; xFactor < 5; xFactor++ {
			for yFactor := 0; yFactor < 5; yFactor++ {
				newPoint := point{
					x: xFactor*size + p.x,
					y: yFactor*size + p.y,
				}
				newRisk := (risk+xFactor+yFactor-1)%9 + 1
				exploded[newPoint] = newRisk
			}
		}
	}
	return exploded
}

func (s *solver) Solve2(input []string) (int, error) {
	points := map[point]int{}
	size := len(input)
	for y, line := range input {
		ints := aoc.IntCharacters(line)
		for x, risk := range ints {
			points[point{x, y}] = risk
		}
	}

	exploded := explode(points, size)
	return dijkstra(exploded, size*5), nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
