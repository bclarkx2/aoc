package main

import (
	"strconv"
	"strings"

	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Dive"
}

func (s *solver) Solve1(input []string) (int, error) {
	x, y := 0, 0

	for _, instruction := range input {
		parts := strings.SplitN(instruction, " ", 2)

		direction := parts[0]
		magnitude, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}

		switch direction {
		case "forward":
			x += magnitude
		case "up":
			y -= magnitude
		case "down":
			y += magnitude
		}
	}

	return x * y, nil
}

func (s *solver) Solve2(input []string) (int, error) {
	aim, x, y := 0, 0, 0

	for _, instruction := range input {
		parts := strings.SplitN(instruction, " ", 2)

		direction := parts[0]
		magnitude, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}

		switch direction {
		case "forward":
			x += magnitude
			y += magnitude * aim
		case "up":
			aim -= magnitude
		case "down":
			aim += magnitude
		}
	}

	return x * y, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
