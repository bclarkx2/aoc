package main

import (
	"flag"
	"log"
	"os"

	"github.com/bclarkx2/aoc"
	"github.com/peterbourgon/ff/v3"
)

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
	if err := ff.Parse(flag.CommandLine, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		log.Fatalf("Error parsing flags: %s", err)
	}

	aoc.Run(*inputFile, &solver{})
}
