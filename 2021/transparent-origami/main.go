package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bclarkx2/aoc"
)

type solver struct{}

func (s *solver) Name() string {
	return "Transparent Origami"
}

type direction string

const (
	vertical   direction = "y"
	horizontal direction = "x"
)

type fold struct {
	direction  direction
	coordinate int
}

type point struct {
	x int
	y int
}

type sheet struct {
	height int
	width  int
	points map[int]map[int]bool
}

func newSheet(points []point) sheet {
	s := sheet{
		points: map[int]map[int]bool{},
	}

	for _, p := range points {
		s.add(p)
	}

	return s
}

func (s *sheet) add(p point) {
	row, ok := s.points[p.y]
	if !ok {
		row = map[int]bool{}
	}
	row[p.x] = true
	s.points[p.y] = row
	if p.y > s.height {
		s.height = p.y
	}
	if p.x > s.width {
		s.width = p.x
	}
}

func (s *sheet) remove(p point) {
	row, ok := s.points[p.y]
	if !ok {
		return
	}

	delete(row, p.x)
	if len(row) == 0 {
		delete(s.points, p.y)
	}
}

func (s *sheet) size() int {
	count := 0
	for _, row := range s.points {
		count += len(row)
	}
	return count
}

func (s *sheet) fold(f fold) {
	switch f.direction {
	case vertical:
		for y, row := range s.points {
			if y < f.coordinate {
				continue
			}
			for x := range row {
				preimage := point{x, y}
				image := point{x, y - 2*(y-f.coordinate)}
				s.remove(preimage)
				s.add(image)
			}
		}
		s.height = f.coordinate - 1
	case horizontal:
		for y, row := range s.points {
			for x := range row {
				if x < f.coordinate {
					continue
				}
				preimage := point{x, y}
				image := point{x - 2*(x-f.coordinate), y}
				s.remove(preimage)
				s.add(image)
			}
		}
		s.width = f.coordinate - 1
	}

}

func (s sheet) String() string {
	points := make([][]string, s.height+1)
	for y := 0; y <= s.height; y++ {
		points[y] = make([]string, s.width+1)
		for x := 0; x <= s.width; x++ {
			if s.points[y][x] {
				points[y][x] = "#"
			} else {
				points[y][x] = "."
			}
		}
	}

	var lines []string
	for _, row := range points {
		line := strings.Join(row, "")
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

var foldRegex = regexp.MustCompile(`fold along (\w{1})=(\d+)`)

func parse(input []string) ([]point, []fold, error) {
	var blankIdx int
	for i, line := range input {
		if line == "" {
			blankIdx = i
		}
	}

	pointLines := input[:blankIdx]
	foldLines := input[blankIdx+1:]

	var points []point
	for _, line := range pointLines {
		coords, err := aoc.Integers(strings.Split(line, ","))
		if err != nil {
			return nil, nil, err
		}
		p := point{
			x: coords[0],
			y: coords[1],
		}
		points = append(points, p)
	}

	var folds []fold
	for _, line := range foldLines {
		matches := foldRegex.FindStringSubmatch(line)
		dir, coordStr := matches[1], matches[2]
		coord, err := strconv.Atoi(coordStr)
		if err != nil {
			return nil, nil, err
		}
		f := fold{
			direction:  direction(dir),
			coordinate: coord,
		}
		folds = append(folds, f)
	}

	return points, folds, nil
}

func (s *solver) Solve1(input []string) (int, error) {
	points, folds, err := parse(input)
	if err != nil {
		return 1, err
	}

	sheet := newSheet(points)
	sheet.fold(folds[0])

	return sheet.size(), nil
}

func (s *solver) Solve2(input []string) (int, error) {
	points, folds, err := parse(input)
	if err != nil {
		return 1, err
	}

	sheet := newSheet(points)
	for _, f := range folds {
		sheet.fold(f)
	}
	fmt.Printf("sheet:\n%s\n", sheet)

	return 0, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
