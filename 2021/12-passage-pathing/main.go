package main

import (
	"fmt"
	"strings"

	"github.com/bclarkx2/aoc"
)

type path struct {
	links          []*cave
	registry       map[string]int
	containsDouble bool
}

func newPath(links ...*cave) path {
	p := path{
		registry: map[string]int{},
	}
	for _, c := range links {
		p.links = append(p.links, c)

		count := p.registry[c.name]
		p.registry[c.name] = count + 1

		if count > 0 && c.size == small {
			p.containsDouble = true
		}
	}
	return p
}

func (p path) last() *cave {
	return p.links[len(p.links)-1]
}

func (p path) extend(end *cave) path {
	var extended []*cave
	for _, c := range p.links {
		extended = append(extended, c)
	}
	extended = append(extended, end)
	return newPath(extended...)
}

func (p path) visits(cave *cave) int {
	return p.registry[cave.name]
}

func (p path) String() string {
	var names []string
	for _, c := range p.links {
		names = append(names, c.name)
	}
	return fmt.Sprintf("%t; %s", p.containsDouble, strings.Join(names, ","))
}

type queue []path

func (q *queue) enqueue(p path) {
	*q = append(*q, p)
}

func (q *queue) dequeue() (path, bool) {
	if len(*q) <= 0 {
		return path{}, false
	}

	p := (*q)[0]
	*q = (*q)[1:]
	return p, true
}

type size int

const (
	small size = iota
	large
)

type cave struct {
	name string
	size size

	neighbors []*cave
}

func (c *cave) addNeighbor(cave *cave) {
	c.neighbors = append(c.neighbors, cave)
}

type caves struct {
	start *cave
	all   map[string]*cave
}

func newCaves(input []string) *caves {
	c := caves{
		all: map[string]*cave{},
	}

	for _, line := range input {
		pieces := strings.Split(line, "-")
		begin, end := pieces[0], pieces[1]
		c.addEdge(begin, end)
	}

	c.start = c.all["start"]
	return &c
}

func (c *caves) upsertNode(name string) *cave {
	node, ok := c.all[name]
	if !ok {
		size := small
		if aoc.IsUpper(name) {
			size = large
		}
		node = &cave{
			name: name,
			size: size,
		}
		c.all[name] = node
	}
	return node
}

func (c *caves) addEdge(begin, end string) {
	beginNode := c.upsertNode(begin)
	endNode := c.upsertNode(end)

	beginNode.addNeighbor(endNode)
	endNode.addNeighbor(beginNode)
}

func (c *caves) paths(doubleLimit int) []path {
	queue := queue{
		newPath(c.start),
	}

	var paths []path
	for {
		// Pop the next path off the queue
		path, ok := queue.dequeue()
		if !ok {
			break
		}

		// Check if this path is finished
		last := path.last()
		if last.name == "end" {
			paths = append(paths, path)
			continue
		}

		// Propose a new path with each neighbor, unless that
		// neighbor is a small cave that's already been
		// visited too many times for that path.
		for _, n := range last.neighbors {
			if n.name == "start" {
				continue
			}
			if n.size == small {
				visits := path.visits(n)
				if path.containsDouble && visits >= 1 {
					continue
				}
				if !path.containsDouble && visits >= doubleLimit {
					continue
				}
			}
			extended := path.extend(n)
			queue.enqueue(extended)
		}
	}

	return paths
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	caves := newCaves(input)
	paths := caves.paths(1)
	return len(paths), nil
}

func (s *solver) Solve2(input []string) (int, error) {
	caves := newCaves(input)
	paths := caves.paths(2)
	return len(paths), nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
