package hexmap

import (
	"fmt"
	"testing"
)

// Tests

func complain(name string) {
	fmt.Println("FAIL ", name)
}

func equal_hex(name string, a Hex, b Hex) bool {
	if !(a.q == b.q && a.s == b.s && a.r == b.r) {
		complain(name)
	}
	return true
}

func equal_offsetcoord(name string, a OffsetCoord, b OffsetCoord) bool {
	if !(a.col == b.col && a.row == b.row) {
		complain(name)
		return false
	}
	return true
}

func equal_doubledcoord(name string, a DoubledCoord, b DoubledCoord) bool {
	if !(a.col == b.col && a.row == b.row) {
		complain(name)
		return false
	}
	return true
}

func equal_int(name string, a int, b int) bool {
	if !(a == b) {
		complain(name)
		return false
	}
	return true
}

func equal_hex_array(name string, a []Hex, b []Hex) bool {
	equal_int(name, len(a), len(b))
	for i := 0; i < len(a); i++ {
		if !equal_hex(name, a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestHexArithmetic(t *testing.T) {
	equal_hex("hex_add", NewHex(4, -10, 6), hex_add(NewHex(1, -3, 2), NewHex(3, -7, 4)))
	equal_hex("hex_subtract", NewHex(-2, 4, -2), hex_subtract(NewHex(1, -3, 2), NewHex(3, -7, 4)))
}

func TestHexDirection(t *testing.T) {
	equal_hex("hex_direction", NewHex(0, -1, 1), hex_direction(2))
}

func TestHexNeighbor(t *testing.T) {
	equal_hex("hex_neighbor", NewHex(1, -3, 2), hex_neighbor(NewHex(1, -2, 1), 2))
}

func TestHexDiagonal(t *testing.T) {
	equal_hex("hex_diagonal", NewHex(-1, -1, 2), hex_diagonal_neighbor(NewHex(1, -2, 1), 3))
}

func TestHexDistance(t *testing.T) {
	equal_int("hex_distance", 7, hex_distance(NewHex(3, -7, 4), NewHex(0, 0, 0)))
}

func TestHexRotateRight(t *testing.T) {
	equal_hex("hex_rotate_right", hex_rotate_right(NewHex(1, -3, 2)), NewHex(3, -2, -1))
}

func TestHexRotateLeft(t *testing.T) {
	equal_hex("hex_rotate_left", hex_rotate_left(NewHex(1, -3, 2)), NewHex(-2, -1, 3))
}

func TestHexRound(t *testing.T) {
	a := NewFractionalHex(0.0, 0.0, 0.0)
	b := NewFractionalHex(1.0, -1.0, 0.0)
	c := NewFractionalHex(0.0, -1.0, 1.0)
	equal_hex("hex_round 1", NewHex(5, -10, 5), hex_round(hex_lerp(NewFractionalHex(0.0, 0.0, 0.0), NewFractionalHex(10.0, -20.0, 10.0), 0.5)))
	equal_hex("hex_round 2", hex_round(a), hex_round(hex_lerp(a, b, 0.499)))
	equal_hex("hex_round 3", hex_round(b), hex_round(hex_lerp(a, b, 0.501)))
	equal_hex("hex_round 4", hex_round(a), hex_round(NewFractionalHex(a.q*0.4+b.q*0.3+c.q*0.3, a.r*0.4+b.r*0.3+c.r*0.3, a.s*0.4+b.s*0.3+c.s*0.3)))
	equal_hex("hex_round 5", hex_round(c), hex_round(NewFractionalHex(a.q*0.3+b.q*0.3+c.q*0.4, a.r*0.3+b.r*0.3+c.r*0.4, a.s*0.3+b.s*0.3+c.s*0.4)))
}

func TestHexLinedraw(t *testing.T) {
	equal_hex_array("hex_linedraw", []Hex{Hex{0, 0, 0}, Hex{0, -1, 1}, Hex{0, -2, 2}, Hex{1, -3, 2}, Hex{1, -4, 3}, Hex{1, -5, 4}}, hex_linedraw(Hex{0, 0, 0}, Hex{1, -5, 4}))
}

func TestLayout(t *testing.T) {
	h := NewHex(3, 4, -7)
	flat := NewLayout(layout_flat, NewPoint(10.0, 15.0), NewPoint(35.0, 71.0))
	equal_hex("layout", h, hex_round(pixel_to_hex(flat, hex_to_pixel(flat, h))))
	pointy := NewLayout(layout_pointy, NewPoint(10.0, 15.0), NewPoint(35.0, 71.0))
	equal_hex("layout", h, hex_round(pixel_to_hex(pointy, hex_to_pixel(pointy, h))))
}

func TestOffsetRoundtrip(t *testing.T) {
	a := NewHex(3, 4, -7)
	b := NewOffsetCoord(1, -3)
	equal_hex("conversion_roundtrip even-q", a, qoffset_to_cube(EVEN, qoffset_from_cube(EVEN, a)))
	equal_offsetcoord("conversion_roundtrip even-q", b, qoffset_from_cube(EVEN, qoffset_to_cube(EVEN, b)))
	equal_hex("conversion_roundtrip odd-q", a, qoffset_to_cube(ODD, qoffset_from_cube(ODD, a)))
	equal_offsetcoord("conversion_roundtrip odd-q", b, qoffset_from_cube(ODD, qoffset_to_cube(ODD, b)))
	equal_hex("conversion_roundtrip even-r", a, roffset_to_cube(EVEN, roffset_from_cube(EVEN, a)))
	equal_offsetcoord("conversion_roundtrip even-r", b, roffset_from_cube(EVEN, roffset_to_cube(EVEN, b)))
	equal_hex("conversion_roundtrip odd-r", a, roffset_to_cube(ODD, roffset_from_cube(ODD, a)))
	equal_offsetcoord("conversion_roundtrip odd-r", b, roffset_from_cube(ODD, roffset_to_cube(ODD, b)))
}

func TestOffsetFromCube(t *testing.T) {
	equal_offsetcoord("offset_from_cube even-q", NewOffsetCoord(1, 3), qoffset_from_cube(EVEN, NewHex(1, 2, -3)))
	equal_offsetcoord("offset_from_cube odd-q", NewOffsetCoord(1, 2), qoffset_from_cube(ODD, NewHex(1, 2, -3)))
}

func TestOffsetToCube(t *testing.T) {
	equal_hex("offset_to_cube even-", NewHex(1, 2, -3), qoffset_to_cube(EVEN, NewOffsetCoord(1, 3)))
	equal_hex("offset_to_cube odd-q", NewHex(1, 2, -3), qoffset_to_cube(ODD, NewOffsetCoord(1, 2)))
}

func TestDoubledRoundtrip(t *testing.T) {
	a := NewHex(3, 4, -7)
	b := NewDoubledCoord(1, -3)
	equal_hex("conversion_roundtrip doubled-q", a, qdoubled_to_cube(qdoubled_from_cube(a)))
	equal_doubledcoord("conversion_roundtrip doubled-q", b, qdoubled_from_cube(qdoubled_to_cube(b)))
	equal_hex("conversion_roundtrip doubled-r", a, rdoubled_to_cube(rdoubled_from_cube(a)))
	equal_doubledcoord("conversion_roundtrip doubled-r", b, rdoubled_from_cube(rdoubled_to_cube(b)))
}

func TestDoubledFromCube(t *testing.T) {
	equal_doubledcoord("doubled_from_cube doubled-q", NewDoubledCoord(1, 5), qdoubled_from_cube(NewHex(1, 2, -3)))
	equal_doubledcoord("doubled_from_cube doubled-r", NewDoubledCoord(4, 2), rdoubled_from_cube(NewHex(1, 2, -3)))
}

func TestDoubledToCube(t *testing.T) {
	equal_hex("doubled_to_cube doubled-q", NewHex(1, 2, -3), qdoubled_to_cube(NewDoubledCoord(1, 5)))
	equal_hex("doubled_to_cube doubled-r", NewHex(1, 2, -3), rdoubled_to_cube(NewDoubledCoord(4, 2)))
}

/*
func TestAll(t *testing.T) {
    TestHexArithmetic(t)
    TestHexDirection(t)
    TestHexNeighbor(t)
    TestHexDiagonal(t)
    TestHexDistance(t)
    TestHexRotateRight(t)
    TestHexRotateLeft(t)
    TestHexRound(t)
    TestHexLinedraw(t)
    TestLayout(t)
    TestOffsetRoundtrip(t)
    TestOffsetFromCube(t)
    TestOffsetToCube(t)
    TestDoubledRoundtrip(t)
    TestDoubledFromCube(t)
    TestDoubledToCube(t)
}
*/
