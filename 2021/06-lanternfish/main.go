package main

import (
	"strings"

	"github.com/bclarkx2/aoc"
)

type key struct {
	age  int
	days int
}

func fish(age int, days int, table map[key]int) int {

	// Check if answer is already contained in the table
	k := key{
		age:  age,
		days: days,
	}
	if answer, ok := table[k]; ok {
		return answer
	}

	// Otherwise, need to compute
	count := 1
	for i := days - age - 1; i >= 0; i -= 7 {
		count += fish(8, i, table)
	}

	// Record answer in the table for later
	table[k] = count

	return count
}

func calculate(input []string, days int) (int, error) {
	state := input[0]

	ageStrs := strings.Split(state, ",")
	ages, err := aoc.Integers(ageStrs)
	if err != nil {
		return 0, err
	}

	count := 0
	table := map[key]int{}
	for _, age := range ages {
		count += fish(age, days, table)
	}

	return count, nil
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	return calculate(input, 80)
}

func (s *solver) Solve2(input []string) (int, error) {
	return calculate(input, 256)
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
