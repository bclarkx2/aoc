package main

import (
	"math"
	"regexp"

	"github.com/bclarkx2/aoc"
)

var promptRegex = regexp.MustCompile(`target area: x=([0-9-]+)..([0-9-]+), y=([0-9-]+)..([0-9-]+)`)

func parse(prompt string) (int, int, int, int, error) {
	matches := promptRegex.FindStringSubmatch(prompt)
	vals, err := aoc.Integers(matches[1:])
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return vals[0], vals[1], vals[2], vals[3], nil
}

func floor(f float64) int {
	return int(math.Floor(f))
}

func ceil(f float64) int {
	return int(math.Ceil(f))
}

func stepsToX(xPosition, xVelocity int) float64 {
	x, vx := float64(xPosition), float64(xVelocity)
	return -0.5*math.Sqrt(4.0*vx*vx+4*vx-8*x+1) + vx + 0.5
}

func velocityForYAfterN(yPosition, steps int) float64 {
	y, n := float64(yPosition), float64(steps)
	return 0.5 * (((2.0 * y) / n) + n - 1.0)
}

func velocityForX(xPosition int) float64 {
	x := float64(xPosition)
	return 0.5*(math.Sqrt(8.0*x)+1.0) - 1.0
}

type velocity struct {
	x int
	y int
}

type nConstraints struct {
	min int
	max int
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	_, _, yMin, _, err := parse(input[0])
	if err != nil {
		return 0, err
	}

	return yMin * (yMin + 1) / 2, nil
}

func (s *solver) Solve2(input []string) (int, error) {
	xMin, xMax, yMin, yMax, err := parse(input[0])
	if err != nil {
		return 0, err
	}

	smallestX := ceil(velocityForX(xMin))

	ns := map[int]nConstraints{}
	for x := smallestX; x <= xMax; x++ {
		nMin := stepsToX(xMin, x)
		nMax := stepsToX(xMax, x)

		if math.IsNaN(nMax) {
			ns[x] = nConstraints{ceil(nMin), math.MaxInt}
		} else if nMin <= nMax {
			ns[x] = nConstraints{ceil(nMin), floor(nMax)}
		}
	}

	valid := map[velocity]bool{}
	for x, nConstraint := range ns {
		for n := nConstraint.min; n <= nConstraint.max; n++ {
			smallestY := ceil(velocityForYAfterN(yMin, n))
			largestY := floor(velocityForYAfterN(yMax, n))

			if smallestY > aoc.Abs(yMin) {
				break
			}

			for y := smallestY; y <= largestY; y++ {
				valid[velocity{x, y}] = true
			}
		}
	}

	return len(valid), nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
