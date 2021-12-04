package aoc

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Solver interface {
	Solve1(input []string) int
	Solve2(input []string) int
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

	fmt.Printf(`
Solving puzzle: %s
Input: %s
Solution 1: %d
Solution 2: %d`,
		solver.Name(),
		inputFile,
		solver.Solve1(lines),
		solver.Solve2(lines),
	)

	return
}

func InputFileFlag() *string {
	var inputFile string
	flag.StringVar(&inputFile, "input", "example.txt", "Input file")
	return &inputFile
}
