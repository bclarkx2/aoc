package aoc

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
	"unicode"

	"github.com/peterbourgon/ff/v3"
)

type Solver interface {
	Solve1(input []string) (int, error)
	Solve2(input []string) (int, error)
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

	e1 := run(lines, solver.Solve1, "Solution 1")
	e2 := run(lines, solver.Solve2, "Solution 2")

	fmt.Printf("\nInput: %s\n", inputFile)
	fmt.Println(e1)
	fmt.Println(e2)
}

func run(input []string, f func([]string) (int, error), label string) execution {
	start := time.Now()
	solution, err := f(input)
	elapsed := time.Since(start)
	return execution{
		label:    label,
		solution: solution,
		err:      err,
		elapsed:  elapsed,
	}
}

type execution struct {
	label    string
	solution int
	err      error
	elapsed  time.Duration
}

func (e execution) String() string {
	var output interface{} = e.solution
	if e.err != nil {
		output = e.err
	}
	return fmt.Sprintf(
		"%s: %v (%vms)",
		e.label,
		output,
		e.elapsed.Milliseconds(),
	)
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

func SortIntsDescending(ints []int) {
	sort.Slice(ints, func(i, j int) bool {
		return ints[i] > ints[j]
	})
}

func Characters(str string) []string {
	var chars []string
	for _, r := range str {
		chars = append(chars, string(r))
	}
	return chars
}

func IntCharacters(str string) []int {
	var ints []int
	for _, r := range str {
		ints = append(ints, int(r-'0'))
	}
	return ints
}

func ContainsStr(lst []string, str string) bool {
	for _, s := range lst {
		if s == str {

			return true
		}
	}
	return false
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
