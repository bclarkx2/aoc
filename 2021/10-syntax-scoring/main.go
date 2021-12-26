package main

import (
	"github.com/bclarkx2/aoc"
)

type stack []string

func (s *stack) push(char string) {
	*s = append(*s, char)
}

func (s *stack) pop() (string, bool) {
	n := len(*s) - 1
	if n < 0 {
		return "", false
	}

	str := (*s)[n]
	*s = (*s)[:n]
	return str, true
}

var (
	openings = []string{"(", "[", "{", "<"}
	values1  = map[string]int{
		")": 3,
		"]": 57,
		"}": 1197,
		">": 25137,
	}
	values2 = map[string]int{
		"(": 1,
		"[": 2,
		"{": 3,
		"<": 4,
	}
)

func matches(opening, closing string) bool {
	switch opening {
	case "(":
		return closing == ")"
	case "[":
		return closing == "]"
	case "{":
		return closing == "}"
	case "<":
		return closing == ">"
	}
	return false
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	score := 0
	for _, line := range input {
		s := stack{}
		for _, char := range aoc.Characters(line) {
			if aoc.ContainsStr(openings, char) {
				s.push(char)
				continue
			}

			corresponding, ok := s.pop()
			if !ok {
				break
			}
			if !matches(corresponding, char) {
				score += values1[char]
				break
			}
		}
	}

	return score, nil
}

func (s *solver) Solve2(input []string) (int, error) {

	var scores []int
lines:
	for _, line := range input {
		s := stack{}
		for _, char := range aoc.Characters(line) {
			if aoc.ContainsStr(openings, char) {
				s.push(char)
				continue
			}

			corresponding, ok := s.pop()
			if !ok {
				s.push(corresponding)
				break
			}
			if !matches(corresponding, char) {
				continue lines
			}
		}

		score := 0
		for opening, ok := s.pop(); ok; opening, ok = s.pop() {
			score = score*5 + values2[opening]
		}
		scores = append(scores, score)
	}

	aoc.SortIntsDescending(scores)
	return scores[len(scores)/2], nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
