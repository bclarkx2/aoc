package main

import (
	"github.com/bclarkx2/aoc"
)

type octopus struct {
	energy int
	x      int
	y      int

	neighbors []*octopus
}

func (o *octopus) neighbor(n *octopus) {
	o.neighbors = append(o.neighbors, n)
}

type octopi [][]*octopus

func (o *octopi) all() []*octopus {
	var a []*octopus
	for _, row := range *o {
		for _, octopus := range row {
			a = append(a, octopus)
		}
	}
	return a
}

func (o *octopi) increment() {
	for _, octopus := range o.all() {
		octopus.energy++
	}
}

func (o *octopi) flash() int {
	flashed := map[*octopus]bool{}

	for {
		numFlashed := 0
		for _, octopus := range o.all() {
			if _, ok := flashed[octopus]; ok {
				continue
			}

			if octopus.energy <= 9 {
				continue
			}

			flashed[octopus] = true
			numFlashed++

			for _, n := range octopus.neighbors {
				n.energy++
			}
		}

		if numFlashed == 0 {
			break
		}
	}

	for o := range flashed {
		o.energy = 0
	}

	return len(flashed)
}

func (o *octopi) size() int {
	return len(o.all())
}

func newOctopi(input []string) octopi {
	// Establish octopi
	var octopi [][]*octopus
	for _, line := range input {
		row := []*octopus{}
		for _, energy := range aoc.IntCharacters(line) {
			octopus := &octopus{
				energy: energy,
			}
			row = append(row, octopus)
		}
		octopi = append(octopi, row)
	}

	// Set relationships
	columns := len(octopi[0])
	rows := len(octopi)
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			octopus := octopi[y][x]
			octopus.x = x
			octopus.y = y

			belowTop := y > 0
			beforeRight := x < columns-1
			aboveBottom := y < rows-1
			pastLeft := x > 0

			if belowTop {
				octopus.neighbor(octopi[y-1][x])
			}
			if beforeRight && belowTop {
				octopus.neighbor(octopi[y-1][x+1])
			}
			if beforeRight {
				octopus.neighbor(octopi[y][x+1])
			}
			if beforeRight && aboveBottom {
				octopus.neighbor(octopi[y+1][x+1])
			}
			if aboveBottom {
				octopus.neighbor(octopi[y+1][x])
			}
			if pastLeft && aboveBottom {
				octopus.neighbor(octopi[y+1][x-1])
			}
			if pastLeft {
				octopus.neighbor(octopi[y][x-1])
			}
			if pastLeft && belowTop {
				octopus.neighbor(octopi[y-1][x-1])
			}
		}
	}

	return octopi
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	octopi := newOctopi(input)

	flashes := 0
	for step := 1; step <= 100; step++ {
		octopi.increment()
		flashes += octopi.flash()
	}

	return flashes, nil
}

func (s *solver) Solve2(input []string) (int, error) {
	octopi := newOctopi(input)

	var step int
	for step = 0; octopi.flash() != octopi.size(); step++ {
		octopi.increment()
	}

	return step, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
