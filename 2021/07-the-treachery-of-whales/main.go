package main

import (
	"strings"

	"github.com/bclarkx2/aoc"
)

type costFunc func(int, int) int

func naive(positions []int, cost costFunc) int {
	min, max := aoc.Min(positions), aoc.Max(positions)

	var minTotal *int
	for target := min; target <= max; target++ {
		total := 0
		for _, p := range positions {
			total += cost(target, p)
		}

		if minTotal == nil || total < *minTotal {
			minTotal = &total
		}
	}

	return *minTotal
}

func calculate(input []string, cost costFunc) (int, error) {
	positionStrs := strings.Split(input[0], ",")
	positions, err := aoc.Integers(positionStrs)
	if err != nil {
		return 0, err
	}

	return naive(positions, cost), nil
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	cost := func(target, position int) int {
		return aoc.AbsDiff(target, position)
	}
	return calculate(input, cost)
}

func (s *solver) Solve2(input []string) (int, error) {
	cost := func(target, position int) int {
		diff := float64(aoc.AbsDiff(target, position))
		raw := ((diff + 1.0) / 2.0) * (0.0 + (diff)*1.0)
		cost := int(raw)
		return cost
	}
	return calculate(input, cost)
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
