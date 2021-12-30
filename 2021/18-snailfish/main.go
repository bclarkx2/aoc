package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/bclarkx2/aoc"
)

type node interface {
	fmt.Stringer
	value() int
	magnitude() int

	explodable(depth int) *pair
	splittable() *leaf

	leftmost() *leaf
	rightmost() *leaf

	setParent(p *pair)
}

type pair struct {
	parent *pair
	left   node
	right  node
}

func newPair(left, right node, parent *pair) *pair {
	p := pair{
		parent: parent,
	}
	p.addChildren(left, right)
	return &p
}

func (p *pair) String() string {
	return fmt.Sprintf("[%s,%s]", p.left, p.right)
}

func (p *pair) value() int {
	panic("value called on branch node")
}

func (p *pair) magnitude() int {
	return 3*p.left.magnitude() + 2*p.right.magnitude()
}

func (p *pair) addChildren(left, right node) {
	p.left, p.right = left, right
	left.setParent(p)
	right.setParent(p)

}

func (p *pair) setParent(parent *pair) {
	p.parent = parent
}

func (p *pair) explode() {
	if peak := p.leftPeak(); peak != nil {
		left := peak.left.rightmost()
		left.val += p.left.value()
	}
	if peak := p.rightPeak(); peak != nil {
		right := peak.right.leftmost()
		right.val += p.right.value()
	}

	l := leaf{
		parent: p.parent,
		val:    0,
	}
	p.parent.replace(p, &l)
}

func (p *pair) replace(old, fresh node) {
	if p.left == old {
		p.left = fresh
	}
	if p.right == old {
		p.right = fresh
	}
}

func (p *pair) leftmost() *leaf {
	return p.left.leftmost()
}

func (p *pair) rightmost() *leaf {
	return p.right.rightmost()
}

func (p *pair) leftPeak() *pair {
	if p.parent == nil {
		return nil
	}
	if p.parent.left == p {
		return p.parent.leftPeak()
	}
	return p.parent
}

func (p *pair) rightPeak() *pair {
	if p.parent == nil {
		return nil
	}
	if p.parent.right == p {
		return p.parent.rightPeak()
	}
	return p.parent
}

func (p *pair) explodable(depth int) *pair {
	_, leftLeaf := p.left.(*leaf)
	_, rightLeaf := p.right.(*leaf)
	if leftLeaf && rightLeaf && depth >= 4 {
		return p
	}
	if e := p.left.explodable(depth + 1); e != nil {
		return e
	}
	if e := p.right.explodable(depth + 1); e != nil {
		return e
	}
	return nil
}

func (p *pair) splittable() *leaf {
	if s := p.left.splittable(); s != nil {
		return s
	}
	if s := p.right.splittable(); s != nil {
		return s
	}
	return nil
}

type leaf struct {
	parent *pair
	val    int
}

func (l *leaf) String() string {
	return fmt.Sprintf("%d", l.val)
}

func (l *leaf) value() int {
	return l.val
}

func (l *leaf) magnitude() int {
	return l.val
}

func (l *leaf) setParent(parent *pair) {
	l.parent = parent
}

func (l *leaf) explodable(depth int) *pair {
	return nil
}

func (l *leaf) splittable() *leaf {
	if l.val >= 10 {
		return l
	}
	return nil
}

func (l *leaf) split() {
	left := leaf{
		val: int(float64(l.val) / 2),
	}
	right := leaf{
		val: int(math.Ceil((float64(l.val) / 2))),
	}
	newNode := newPair(&left, &right, l.parent)

	l.parent.replace(l, newNode)
}

func (l *leaf) leftmost() *leaf {
	return l
}

func (l *leaf) rightmost() *leaf {
	return l
}

type number struct {
	root node
}

func (n number) String() string {
	return n.root.String()
}

func (n *number) Reduce() {
	for {
		if e := n.root.explodable(0); e != nil {
			e.explode()
		} else if s := n.root.splittable(); s != nil {
			s.split()
		} else {
			break
		}
	}
}

func (n *number) Add(operand number) {
	if n.root == nil {
		n.root = operand.root
		return
	}
	n.root = newPair(n.root, operand.root, nil)
	n.Reduce()
}

func (n *number) Magnitude() int {
	return n.root.magnitude()
}

func newNumber(str string) (number, error) {
	ns := stack{}
	for _, char := range aoc.Characters(str) {
		switch char {
		case ",":
			break
		case "[":
			node := pair{}
			ns.push(&node)
		case "]":
			nodes, ok := ns.popN(3)
			if !ok {
				return number{}, errors.New("missing nodes")
			}
			right, left, parent := nodes[0], nodes[1], nodes[2].(*pair)
			parent.addChildren(left, right)
			ns.push(parent)
		default:
			value, err := strconv.Atoi(char)
			if err != nil {
				return number{}, err
			}
			leaf := leaf{
				val: value,
			}
			ns.push(&leaf)
		}
	}

	root, ok := ns.pop()
	if !ok {
		return number{}, errors.New("missing root")
	}

	return number{
		root: root,
	}, nil
}

type stack []node

func (s *stack) push(n node) {
	*s = append(*s, n)
}

func (s *stack) pop() (node, bool) {
	n := len(*s) - 1
	if n < 0 {
		return nil, false
	}

	node := (*s)[n]
	*s = (*s)[:n]
	return node, true
}

func (s *stack) popN(n int) ([]node, bool) {
	nodes := make([]node, n)
	ok := false
	for i := 0; i < n; i++ {
		nodes[i], ok = s.pop()
		if !ok {
			return nil, false
		}
	}
	return nodes, true
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	var sum number
	for _, line := range input {
		n, err := newNumber(line)
		if err != nil {
			return 0, err
		}

		sum.Add(n)
	}

	return sum.Magnitude(), nil
}

func (s *solver) Solve2(input []string) (int, error) {
	var max int
	for _, l1 := range input {
		for _, l2 := range input {
			n1, err := newNumber(l1)
			if err != nil {
				return 0, err
			}

			n2, err := newNumber(l2)
			if err != nil {
				return 0, err
			}

			n1.Add(n2)
			mag := n1.Magnitude()
			max = aoc.Max([]int{mag, max})
		}
	}

	return max, nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
