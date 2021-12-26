package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bclarkx2/aoc"
)

type note struct {
	inputs  []string
	outputs []string
}

func parse(line string) note {
	pieces := strings.Split(line, " | ")
	return note{
		inputs:  strings.Split(pieces[0], " "),
		outputs: strings.Split(pieces[1], " "),
	}
}

var sevenSegmentValues = map[string]int{
	"abcefg":  0,
	"cf":      1,
	"acdeg":   2,
	"acdfg":   3,
	"bcdf":    4,
	"abdfg":   5,
	"abdefg":  6,
	"acf":     7,
	"abcdefg": 8,
	"abcdfg":  9,
}

func update(counts map[string]int, digit string) {
	for _, segRune := range digit {
		seg := string(segRune)
		counts[seg] = counts[seg] + 1
	}
}

func segments(digit string) map[string]bool {
	m := map[string]bool{}
	for _, segRune := range digit {
		seg := string(segRune)
		m[seg] = true
	}
	return m
}

func count(digit string, subset map[string]bool) int {
	c := 0
	for _, segRune := range digit {
		seg := string(segRune)
		if _, ok := subset[seg]; ok {
			c++
		}
	}
	return c
}

func value(digit string, mappings map[string]string) int {
	var mapped []string
	for _, segRune := range digit {
		seg := string(segRune)
		mapped = append(mapped, mappings[seg])
	}

	sort.Strings(mapped)
	mappedDigit := strings.Join(mapped, "")

	return sevenSegmentValues[mappedDigit]
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	var notes []note
	for _, line := range input {
		notes = append(notes, parse(line))
	}

	count := 0
	for _, note := range notes {
		for _, digit := range note.outputs {
			l := len(digit)
			if l == 2 || l == 3 || l == 4 || l == 7 {
				count++
			}
		}
	}

	return count, nil
}

func (s *solver) Solve2(input []string) (int, error) {
	var notes []note
	for _, line := range input {
		notes = append(notes, parse(line))
	}

	sum := 0
	for _, note := range notes {

		digits := map[string]int{}
		segs := map[int]map[string]bool{}
		fives := []string{}
		sixes := []string{}

		counts := map[string]int{}

		for _, digit := range note.inputs {
			switch len(digit) {
			case 2:
				digits[digit] = 1
				segs[1] = segments(digit)
				update(counts, digit)
			case 3:
				digits[digit] = 7
				segs[7] = segments(digit)
				update(counts, digit)
			case 4:
				digits[digit] = 4
				segs[4] = segments(digit)
				update(counts, digit)
			case 5:
				fives = append(fives, digit)
			case 6:
				sixes = append(sixes, digit)
			case 7:
				digits[digit] = 8
				segs[8] = segments(digit)
				update(counts, digit)
			}

		}

		// counts
		// 4: c, f
		// 2: a, b, d
		// 1: e, g
		cOrF := map[string]bool{}
		aOrBOrD := map[string]bool{}
		eOrG := map[string]bool{}
		for seg, count := range counts {
			switch count {
			case 1:
				eOrG[seg] = true
			case 2:
				aOrBOrD[seg] = true
			case 4:
				cOrF[seg] = true
			}
		}

		// Hold all the original -> scrambled segment mappings
		mapping := map[string]string{}

		// a is the segment of input 7 that is not
		// a segment in input 1
		for seg := range segs[7] {
			if _, ok := segs[1][seg]; !ok {
				mapping["a"] = seg
			}
		}

		bOrD := map[string]bool{}
		for seg := range aOrBOrD {
			if seg == mapping["a"] {
				continue
			}
			bOrD[seg] = true
		}

		// input 9 is the input with six segments that only
		// contains either e or g
		for _, digit := range sixes {
			if count(digit, eOrG) == 1 {
				segs[9] = segments(digit)
				digits[digit] = 9
			}
		}

		// the segment that isn't c, f, b, d, or a in
		// the nine input must be g
		for seg := range segs[9] {
			_, cf := cOrF[seg]
			_, bd := bOrD[seg]
			a := (seg == mapping["a"])
			if !cf && !bd && !a {
				mapping["g"] = seg
			}
		}

		// e must be whichever is not g from eOrG
		for seg := range eOrG {
			if seg != mapping["g"] {
				mapping["e"] = seg
			}
		}

		// identify 0 and 6
		for _, digit := range sixes {
			is9 := (digits[digit] == 9)
			bdCount := count(digit, bOrD)
			if !is9 && bdCount == 1 {
				digits[digit] = 0
				segs[0] = segments(digit)
			} else if !is9 && bdCount == 2 {
				digits[digit] = 6
				segs[6] = segments(digit)
			}
		}

		// b is the one that's in the zero input
		// that isn't c, f, a, e, or g
		for segRune := range segs[0] {
			seg := string(segRune)
			isA := (seg == mapping["a"])
			_, isCOrF := cOrF[seg]
			_, isEOrG := eOrG[seg]
			if !isA && !isCOrF && !isEOrG {
				mapping["b"] = seg
			}
		}

		// d is the member of bOrD that isn't b
		for seg := range bOrD {
			if seg != mapping["b"] {
				mapping["d"] = seg
			}
		}

		// f is the member of input 6 that isn't
		// a, b, c, d, e, or g
		for segRune := range segs[6] {
			seg := string(segRune)
			_, isABD := aOrBOrD[seg]
			_, isEG := eOrG[seg]
			_, isCF := cOrF[seg]
			if !isABD && !isEG && isCF {
				mapping["f"] = seg
			}
		}

		// c is the member of cOrF that isn't f
		for seg := range cOrF {
			if seg != mapping["f"] {
				mapping["c"] = seg
			}
		}

		// invert the mapping to calculate the output value
		scrambledToOriginal := map[string]string{}
		for k, v := range mapping {
			scrambledToOriginal[v] = k
		}

		// find value of each output digit
		var values []int
		for _, output := range note.outputs {
			values = append(values, value(output, scrambledToOriginal))
		}

		var strs []string
		for _, v := range values {
			strs = append(strs, fmt.Sprintf("%d", v))
		}
		outputStr := strings.Join(strs, "")
		output, err := strconv.Atoi(outputStr)
		if err != nil {
			return 0, err
		}

		sum += output
	}

	return sum, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
