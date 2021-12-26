package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bclarkx2/aoc"
)

type opportunity map[int]bool

func newOpportunity(vals []int) opportunity {
	opp := opportunity{}
	for _, val := range vals {
		opp[val] = true
	}
	return opp
}

type board struct {
	id      int
	numbers map[int]bool
	opps    map[int][]opportunity
}

func newBoard(vals [][]int, id int) board {

	opps := map[int][]opportunity{}
	numbers := map[int]bool{}

	// process rows
	for _, row := range vals {
		// record each horizontal opp
		opp := newOpportunity(row)

		for _, num := range row {
			// register the opp for that number
			opps[num] = append(opps[num], opp)

			// also record the number overall
			numbers[num] = true
		}
	}

	// process columns
	for j := 0; j < len(vals[0]); j++ {
		var column []int
		for i := 0; i < len(vals); i++ {
			column = append(column, vals[i][j])
		}

		opp := newOpportunity(column)
		for _, val := range column {
			opps[val] = append(opps[val], opp)
		}
	}

	return board{
		id:      id,
		numbers: numbers,
		opps:    opps,
	}
}

func (b *board) record(num int) bool {
	// remove number from unmarked list
	delete(b.numbers, num)

	// cross off number from opps
	for _, opp := range b.opps[num] {
		delete(opp, num)
		if len(opp) == 0 {
			return true
		}
	}

	return false
}

func (b *board) unmarkedSum() int {
	sum := 0
	for num := range b.numbers {
		sum += num
	}
	return sum
}

func parse(lines []string) ([]int, map[int]board, error) {
	drawStrs := lines[0]
	remaining := lines[2:]

	var boardLines []string
	for _, line := range remaining {
		if line != "" {
			boardLines = append(boardLines, line)
		}
	}

	boards := map[int]board{}
	id := 0
	for start := 0; start < len(boardLines); start += 5 {
		strLines := boardLines[start : start+5]

		intLines := make([][]int, 5)
		for row, strLine := range strLines {
			numberStrs := strings.Fields(strLine)
			for _, numberStr := range numberStrs {
				val, err := strconv.Atoi(numberStr)
				if err != nil {
					return nil, nil, err
				}

				intLines[row] = append(intLines[row], val)
			}
		}
		board := newBoard(intLines, id)
		id++
		boards[id] = board

	}

	var draws []int
	for _, drawStr := range strings.Split(drawStrs, ",") {
		draw, err := strconv.Atoi(drawStr)
		if err != nil {
			return nil, nil, err
		}
		draws = append(draws, draw)
	}

	return draws, boards, nil
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	draws, boards, err := parse(input)
	if err != nil {
		return 0, err
	}

	for _, draw := range draws {
		for _, board := range boards {
			if won := board.record(draw); won {
				return board.unmarkedSum() * draw, nil
			}
		}
	}

	return 0, errors.New("no answer found")
}

func (s *solver) Solve2(input []string) (int, error) {
	draws, boards, err := parse(input)
	if err != nil {
		return 0, err
	}

	winners := map[int]bool{}
	for _, draw := range draws {
		for id, board := range boards {
			// can ignore boards that have already won
			if _, won := winners[id]; won {
				continue
			}

			// can keep going if this board doesn't win
			if won := board.record(draw); !won {
				continue
			}

			winners[id] = true
			if len(winners) == len(boards) {
				return board.unmarkedSum() * draw, nil
			}
		}
	}

	return 0, errors.New("no answer found")
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
