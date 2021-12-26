package main

import (
	"regexp"

	"github.com/bclarkx2/aoc"
)

var lineExp = regexp.MustCompile(`(\d+),(\d+) -> (\d+),(\d+)`)

type point struct {
	x int
	y int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func midpoints(begin, end point) []point {
	stepsX := max(begin.x, end.x) - min(begin.x, end.x)
	stepsY := max(begin.y, end.y) - min(begin.y, end.y)
	steps := max(stepsX, stepsY)

	deltaX, deltaY := 0, 0

	if stepsX > 0 {
		deltaX = (end.x - begin.x) / stepsX
	}

	if stepsY > 0 {
		deltaY = (end.y - begin.y) / stepsY
	}

	var addrs []point
	for step := 0; step <= steps; step++ {
		addr := point{
			x: begin.x + step*deltaX,
			y: begin.y + step*deltaY,
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {

	heatmap := map[point]int{}
	doubles := 0
	for _, line := range input {
		matches := lineExp.FindStringSubmatch(line)[1:]
		nums, err := aoc.Integers(matches)
		if err != nil {
			return 0, err
		}

		begin := point{nums[0], nums[1]}
		end := point{nums[2], nums[3]}

		var addrs []point
		if begin.x == end.x {
			for y := min(begin.y, end.y); y <= max(begin.y, end.y); y++ {
				addr := point{begin.x, y}
				addrs = append(addrs, addr)
			}
		} else if begin.y == end.y {
			for x := min(begin.x, end.x); x <= max(begin.x, end.x); x++ {
				addr := point{x, begin.y}
				addrs = append(addrs, addr)
			}
		}

		for _, addr := range addrs {
			heatmap[addr] += 1
			if heatmap[addr] == 2 {
				doubles++
			}
		}
	}

	return doubles, nil
}

func (s *solver) Solve2(input []string) (int, error) {

	heatmap := map[point]int{}
	doubles := 0
	for _, line := range input {
		matches := lineExp.FindStringSubmatch(line)[1:]
		nums, err := aoc.Integers(matches)
		if err != nil {
			return 0, err
		}

		begin := point{nums[0], nums[1]}
		end := point{nums[2], nums[3]}

		points := midpoints(begin, end)
		for _, point := range points {
			heatmap[point] += 1
			if heatmap[point] == 2 {
				doubles++
			}
		}
	}

	return doubles, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
