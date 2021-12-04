package main

import "bclarkx2/aoc"

type solver struct{}

func (s *solver) Name() string {
	return "Sonar Sweep"
}

func (s *solver) Solve1(input []string) int {
	return 42
}

func (s *solver) Solve2(input []string) int {
	return 442200
}

var (
	inputFile = aoc.InputFileFlag()
)

func main() {
	aoc.Run(*inputFile, &solver{})
}
