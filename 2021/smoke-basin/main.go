package main

import (
	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Smoke Basin"
}

type point struct {
	height int
	x      int
	y      int

	above *point
	right *point
	below *point
	left  *point
}

func newPoints(input []string) [][]*point {
	// Establish points
	var points [][]*point
	for _, line := range input {
		row := []*point{}
		for _, heightRune := range line {
			point := &point{
				height: int(heightRune - '0'),
			}
			row = append(row, point)
		}
		points = append(points, row)
	}

	// Set relationships
	columns := len(points[0])
	rows := len(points)
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			point := points[y][x]
			point.x = x
			point.y = y
			if y > 0 {
				point.above = points[y-1][x]
			}
			if x < columns-1 {
				point.right = points[y][x+1]
			}
			if y < rows-1 {
				point.below = points[y+1][x]
			}
			if x > 0 {
				point.left = points[y][x-1]
			}
		}
	}

	return points
}

func (p *point) Risk() int {
	return p.height + 1
}

func (p *point) IsLowPoint() bool {
	if p.above != nil && p.above.height <= p.height {
		return false
	}
	if p.right != nil && p.right.height <= p.height {
		return false
	}
	if p.below != nil && p.below.height <= p.height {
		return false
	}
	if p.left != nil && p.left.height <= p.height {
		return false
	}
	return true
}

func (p *point) BasinSize() int {
	seen := map[*point]bool{}
	p.basin(seen)
	return len(seen)
}

func (p *point) basin(seen map[*point]bool) {
	// if this is a peak, it does not increase the basin
	if p.height == 9 {
		return
	}

	// check if this point is already in the basin
	if _, ok := seen[p]; ok {
		return
	}

	// otherwise, always add this point to the basin
	seen[p] = true

	// add neighbor's basins if they are uphill
	if p.above != nil && p.above.height >= p.height {
		p.above.basin(seen)
	}
	if p.right != nil && p.right.height >= p.height {
		p.right.basin(seen)
	}
	if p.below != nil && p.below.height >= p.height {
		p.below.basin(seen)
	}
	if p.left != nil && p.left.height >= p.height {
		p.left.basin(seen)
	}

	return
}

func (s *solver) Solve1(input []string) (int, error) {
	points := newPoints(input)

	risk := 0
	for _, row := range points {
		for _, point := range row {
			if point.IsLowPoint() {
				risk += point.Risk()
			}
		}
	}

	return risk, nil
}

func (s *solver) Solve2(input []string) (int, error) {
	points := newPoints(input)

	var sizes []int
	for _, row := range points {
		for _, point := range row {
			if point.IsLowPoint() {
				sizes = append(sizes, point.BasinSize())
			}
		}
	}

	aoc.SortIntsDescending(sizes)
	return sizes[0] * sizes[1] * sizes[2], nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
