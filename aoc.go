package aoc

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/peterbourgon/ff/v3"
)

type Solver interface {
	Solve1(input []string) (int, error)
	Solve2(input []string) (int, error)
	Name() string
}

func Run(inputFile string, solver Solver) {
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s", err)
		return
	}

	var output1, output2 string

	solution1, err1 := solver.Solve1(lines)
	if err1 != nil {
		output1 = err1.Error()
	} else {
		output1 = fmt.Sprintf("%d", solution1)
	}

	solution2, err2 := solver.Solve2(lines)
	if err2 != nil {
		output2 = err2.Error()
	} else {
		output2 = fmt.Sprintf("%d", solution2)
	}

	fmt.Printf(`
%s
Input: %s
Solution 1: %s
Solution 2: %s`,
		solver.Name(),
		inputFile,
		output1,
		output2,
	)

	return
}

func ParseInputFile() string {
	var inputFile string
	flag.StringVar(&inputFile, "input", "puzzle.txt", "Input file")

	if err := ff.Parse(flag.CommandLine, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		log.Fatalf("Error parsing flags: %s", err)
	}

	return inputFile
}

func Integers(strs []string) ([]int, error) {
	var integers []int
	for _, str := range strs {
		integer, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		integers = append(integers, integer)
	}
	return integers, nil
}

// Integer power: compute a**b using binary powering algorithm
// See Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
func Pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func Abs(n int) int {
	return AbsDiff(n, 0)
}

func AbsDiff(n, m int) int {
	if n < m {
		return m - n
	}
	return n - m
}

func Min(ints []int) int {
	if len(ints) < 1 {
		panic("Min called on zero length array")
	}

	min := ints[0]
	for _, i := range ints {
		if i < min {
			min = i
		}
	}

	return min
}

func Max(ints []int) int {
	if len(ints) < 1 {
		panic("Max called on zero length array")
	}

	min := ints[0]
	for _, i := range ints {
		if i > min {
			min = i
		}
	}

	return min
}
