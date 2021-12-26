package main

import (
	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	return 1, nil
}

func (s *solver) Solve2(input []string) (int, error) {
	return 2, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
