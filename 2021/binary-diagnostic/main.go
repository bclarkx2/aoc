package main

import (
	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Binary-Diagnostics"
}

func (s *solver) Solve1(input []string) (int, error) {
	sums := make([]int, len(input[0]))
	for _, i := range input {
		for place, d := range i {
			sums[place] += int(d - '0')
		}
	}

	gamma, epsilon := 0, 0
	for idx, sum := range sums {
		weight := len(sums) - idx - 1
		increase := aoc.Pow(2, weight)
		if sum < (len(input) / 2) {
			epsilon += increase
		} else {
			gamma += increase
		}
	}

	return epsilon * gamma, nil
}

func find(input []string, cmp func(zeroes, ones []string) []string) string {
	candidates := map[string]bool{}
	for _, i := range input {
		candidates[i] = true
	}

	for place := range input[0] {
		var zeroes, ones []string
		for str := range candidates {
			if val := int(str[place]) - '0'; val == 0 {
				zeroes = append(zeroes, str)
			} else {
				ones = append(ones, str)
			}
		}

		for _, str := range cmp(zeroes, ones) {
			delete(candidates, str)
		}

		if len(candidates) == 1 {
			for value := range candidates {
				return value
			}
		}
	}

	return ""
}

func (s *solver) Solve2(input []string) (int, error) {
	mostStr := find(input, func(zeroes, ones []string) []string {
		if len(ones) >= len(zeroes) {
			return zeroes
		}
		return ones
	})
	leastStr := find(input, func(zeroes, ones []string) []string {
		if len(ones) >= len(zeroes) {
			return ones
		}
		return zeroes
	})

	most, least := 0, 0
	for place := range input[0] {
		weight := len(input[0]) - place - 1
		increase := aoc.Pow(2, weight)
		if int(mostStr[place])-'0' == 1 {
			most += increase
		}
		if int(leastStr[place])-'0' == 1 {
			least += increase
		}
	}

	return most * least, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
