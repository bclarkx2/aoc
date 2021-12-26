package main

import (
	"math"

	"github.com/bclarkx2/aoc"
)

const (
	typeSum         = 0
	typeProduct     = 1
	typeMinimum     = 2
	typeMaximum     = 3
	typeLiteral     = 4
	typeGreaterThan = 5
	typeLessThan    = 6
	typeEqualTo     = 7

	lengthTypeBits   = 0
	lengthTypeNumber = 1
)

type packet interface {
	Version() int
	Value() int
}

type operator struct {
	version int
	typeID  int

	packets []packet
}

func (o operator) Version() int {
	sum := o.version
	for _, p := range o.packets {
		sum += p.Version()
	}
	return sum
}

func (o operator) Value() int {
	value := 0
	switch o.typeID {
	case typeSum:
		for _, p := range o.packets {
			value += p.Value()
		}
	case typeProduct:
		value = 1
		for _, p := range o.packets {
			value *= p.Value()
		}
	case typeMinimum:
		value = math.MaxInt
		for _, p := range o.packets {
			if subVal := p.Value(); subVal < value {
				value = subVal
			}
		}
	case typeMaximum:
		for _, p := range o.packets {
			if subVal := p.Value(); subVal > value {
				value = subVal
			}
		}
	case typeGreaterThan:
		first, second := o.packets[0], o.packets[1]
		if first.Value() > second.Value() {
			value = 1
		}
	case typeLessThan:
		first, second := o.packets[0], o.packets[1]
		if first.Value() < second.Value() {
			value = 1
		}
	case typeEqualTo:
		first, second := o.packets[0], o.packets[1]
		if first.Value() == second.Value() {
			value = 1
		}
	}
	return value
}

type literal struct {
	version int
	typeID  int

	value int
}

func (l literal) Version() int {
	return l.version
}

func (l literal) Value() int {
	return l.value
}

type bits []uint

func toBits(hex string) bits {
	bits := make([]uint, len(hex)*4)
	for i, char := range hex {
		addr := i * 4
		switch char {
		case '0':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 0, 0, 0
		case '1':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 0, 0, 1
		case '2':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 0, 1, 0
		case '3':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 0, 1, 1
		case '4':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 1, 0, 0
		case '5':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 1, 0, 1
		case '6':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 1, 1, 0
		case '7':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 0, 1, 1, 1
		case '8':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 0, 0, 0
		case '9':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 0, 0, 1
		case 'A':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 0, 1, 0
		case 'B':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 0, 1, 1
		case 'C':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 1, 0, 0
		case 'D':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 1, 0, 1
		case 'E':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 1, 1, 0
		case 'F':
			bits[addr], bits[addr+1], bits[addr+2], bits[addr+3] = 1, 1, 1, 1
		}
	}
	return bits
}

func toNum(bits bits) int {
	digits := len(bits)
	total := 0
	for place := 0; place < digits; place++ {
		if bits[place] == 1 {
			value := digits - place - 1
			total += aoc.Pow(2, value)
		}
	}
	return total
}

type cursor struct {
	bits bits
	i    int
}

func newCursor(bits bits) *cursor {
	return &cursor{
		bits: bits,
		i:    0,
	}
}

func (c *cursor) position() int {
	return c.i
}

func (c *cursor) peek(digits int) int {
	return toNum(c.bits[c.i : c.i+digits])
}

func (c *cursor) pop(digits int) int {
	val := c.peek(digits)
	c.i += digits
	return val
}

func (c *cursor) popBits(digits int) bits {
	bits := c.bits[c.i : c.i+digits]
	c.i += digits
	return bits
}

func parseLiteral(c *cursor, version, typeID int) literal {
	bits := bits{}
	for {
		continuation := c.pop(1)
		chunk := c.popBits(4)
		bits = append(bits, chunk...)

		if continuation == 0 {
			break
		}
	}

	return literal{
		version: version,
		typeID:  typeID,
		value:   toNum(bits),
	}
}

func parseSubpacketsByLength(c *cursor, length int) []packet {
	begin := c.position()

	var packets []packet
	for c.position() < begin+length {
		packets = append(packets, parse(c))
	}

	return packets
}

func parseSubpacketsByNumber(c *cursor, num int) []packet {
	packets := make([]packet, num)
	for i := 0; i < num; i++ {
		packets[i] = parse(c)
	}
	return packets
}

func parseOperator(c *cursor, version, typeID int) operator {
	lengthTypeID := c.pop(1)

	var packets []packet
	switch lengthTypeID {
	case lengthTypeBits:
		length := c.pop(15)
		packets = parseSubpacketsByLength(c, length)
	case lengthTypeNumber:
		num := c.pop(11)
		packets = parseSubpacketsByNumber(c, num)
	}

	return operator{
		version: version,
		typeID:  typeID,
		packets: packets,
	}
}

func parse(c *cursor) packet {
	version := c.pop(3)
	typeID := c.pop(3)

	switch typeID {
	case typeLiteral:
		return parseLiteral(c, version, typeID)
	default:
		return parseOperator(c, version, typeID)
	}
}

type solver struct{}

func (s *solver) Solve1(input []string) (int, error) {
	bits := toBits(input[0])
	c := newCursor(bits)
	packet := parse(c)
	return packet.Version(), nil
}

func (s *solver) Solve2(input []string) (int, error) {
	bits := toBits(input[0])
	c := newCursor(bits)
	packet := parse(c)
	return packet.Value(), nil
}

func main() {
	aoc.Run(aoc.ParseInputFile(), &solver{})
}
