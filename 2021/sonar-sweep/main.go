package main

import (
	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Sonar Sweep"
}

func (s *solver) Solve1(input []string) (int, error) {
	depths, err := aoc.Integers(input)
	if err != nil {
		return 0, err
	}

	return increases(depths), nil
}

func (s *solver) Solve2(input []string) (int, error) {
	depths, err := aoc.Integers(input)
	if err != nil {
		return 0, err
	}

	windows := make([]int, len(depths)-2)
	for i := 0; i < len(depths)-2; i++ {
		windows[i] = depths[i] + depths[i+1] + depths[i+2]
	}

	return increases(windows), nil
}

func increases(depths []int) int {
	count := 0
	for i := 0; i < len(depths)-1; i++ {
		if depths[i] < depths[i+1] {
			count++
		}
	}

	return count
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
