package main

import (
	"strings"

	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Extended Polymerization"
}

type pair struct {
	first  string
	second string
}

type rule struct {
	pair   pair
	result string
}

type count map[string]int

func newCount(elems ...string) count {
	c := count{}
	for _, elem := range elems {
		c[elem] = c[elem] + 1
	}
	return c
}

func merge(c1, c2 count, intersect string) count {
	c := count{}
	for elem, num := range c1 {
		c[elem] = c[elem] + num
	}
	for elem, num := range c2 {
		c[elem] = c[elem] + num
	}
	c[intersect] = c[intersect] - 1
	return c
}

func solve(chain string, ruleList []rule, n int) int {
	rules := map[pair]string{}
	for _, rule := range ruleList {
		rules[rule.pair] = rule.result
	}

	characters := map[string]bool{}
	for _, char := range aoc.Characters(chain) {
		characters[char] = true
	}
	for _, rule := range ruleList {
		characters[rule.result] = true
	}

	table := map[pair][]count{}
	for first := range characters {
		for second := range characters {
			row := make([]count, n+1)
			pair := pair{
				first:  first,
				second: second,
			}

			row[0] = newCount(first, second)
			table[pair] = row
		}
	}

	for step := 1; step <= n; step++ {
		for p, row := range table {
			if result, ok := rules[p]; ok {
				pair1 := pair{
					first:  p.first,
					second: result,
				}
				pair2 := pair{
					first:  result,
					second: p.second,
				}
				row[step] = merge(table[pair1][step-1], table[pair2][step-1], result)
			} else {
				row[step] = row[step-1]
			}
		}
	}

	chars := aoc.Characters(chain)
	finalCount := newCount(chars[0])
	for i := 0; i < len(chars)-1; i++ {
		pair := pair{
			first:  chars[i],
			second: chars[i+1],
		}
		finalCount = merge(finalCount, table[pair][n], chars[i])
	}

	min, max := chars[0], chars[0]
	for char, val := range finalCount {
		if val < finalCount[min] {
			min = char
		}
		if val > finalCount[max] {
			max = char
		}
	}

	return finalCount[max] - finalCount[min]
}

func parse(input []string) (string, []rule) {
	chain := input[0]
	ruleStrs := input[2:]

	var rules []rule
	for _, line := range ruleStrs {
		pieces := strings.Split(line, " -> ")
		pairStr, result := pieces[0], pieces[1]
		first, second := string(pairStr[0]), string(pairStr[1])
		r := rule{
			pair: pair{
				first:  first,
				second: second,
			},
			result: result,
		}
		rules = append(rules, r)
	}

	return chain, rules
}
func (s *solver) Solve1(input []string) (int, error) {
	chain, rules := parse(input)
	return solve(chain, rules, 10), nil
}

func (s *solver) Solve2(input []string) (int, error) {
	chain, rules := parse(input)
	return solve(chain, rules, 40), nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
