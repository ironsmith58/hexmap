// hexmap test, use 'go test'
package hexmap

import (
	"fmt"
	"testing"
)

// Hexmap Tests

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
	equal_hex("hex_add", NewHex(4, -10, 6), HexAdd(NewHex(1, -3, 2), NewHex(3, -7, 4)))
	equal_hex("hex_subtract", NewHex(-2, 4, -2), HexSubtract(NewHex(1, -3, 2), NewHex(3, -7, 4)))
}

func TestHexDirection(t *testing.T) {
	equal_hex("hex_direction", NewHex(0, -1, 1), HexDirection(2))
}

func TestHexNeighbor(t *testing.T) {
	equal_hex("hex_neighbor", NewHex(1, -3, 2), HexNeighbor(NewHex(1, -2, 1), 2))
}

func TestHexDiagonal(t *testing.T) {
	equal_hex("hex_diagonal", NewHex(-1, -1, 2), HexDiagonalNeighbor(NewHex(1, -2, 1), 3))
}

func TestHexDistance(t *testing.T) {
	equal_int("hex_distance", 7, HexDistance(NewHex(3, -7, 4), NewHex(0, 0, 0)))
}

func TestHexRotateRight(t *testing.T) {
	equal_hex("hex_rotate_right", HexRotateRight(NewHex(1, -3, 2)), NewHex(3, -2, -1))
}

func TestHexRotateLeft(t *testing.T) {
	equal_hex("hex_rotate_left", HexRotateLeft(NewHex(1, -3, 2)), NewHex(-2, -1, 3))
}

func TestHexRound(t *testing.T) {
	a := NewFractionalHex(0.0, 0.0, 0.0)
	b := NewFractionalHex(1.0, -1.0, 0.0)
	c := NewFractionalHex(0.0, -1.0, 1.0)
	equal_hex("hex_round 1", NewHex(5, -10, 5), HexRound(HexLerp(NewFractionalHex(0.0, 0.0, 0.0), NewFractionalHex(10.0, -20.0, 10.0), 0.5)))
	equal_hex("hex_round 2", HexRound(a), HexRound(HexLerp(a, b, 0.499)))
	equal_hex("hex_round 3", HexRound(b), HexRound(HexLerp(a, b, 0.501)))
	equal_hex("hex_round 4", HexRound(a), HexRound(NewFractionalHex(a.q*0.4+b.q*0.3+c.q*0.3, a.r*0.4+b.r*0.3+c.r*0.3, a.s*0.4+b.s*0.3+c.s*0.3)))
	equal_hex("hex_round 5", HexRound(c), HexRound(NewFractionalHex(a.q*0.3+b.q*0.3+c.q*0.4, a.r*0.3+b.r*0.3+c.r*0.4, a.s*0.3+b.s*0.3+c.s*0.4)))
}

func TestHexLinedraw(t *testing.T) {
	equal_hex_array("hex_linedraw", []Hex{Hex{0, 0, 0}, Hex{0, -1, 1}, Hex{0, -2, 2}, Hex{1, -3, 2}, Hex{1, -4, 3}, Hex{1, -5, 4}}, HexLineDraw(Hex{0, 0, 0}, Hex{1, -5, 4}))
}

func TestLayout(t *testing.T) {
	h := NewHex(3, 4, -7)
	flat := NewLayout(layoutFlat, NewPoint(10.0, 15.0), NewPoint(35.0, 71.0))
	equal_hex("layout", h, HexRound(PixelToHex(flat, HexToPixel(flat, h))))
	pointy := NewLayout(layoutPointy, NewPoint(10.0, 15.0), NewPoint(35.0, 71.0))
	equal_hex("layout", h, HexRound(PixelToHex(pointy, HexToPixel(pointy, h))))
}

func TestOffsetRoundtrip(t *testing.T) {
	a := NewHex(3, 4, -7)
	b := NewOffsetCoord(1, -3)
	equal_hex("conversion_roundtrip even-q", a, QoffsetToCube(EVEN, QoffsetFromCube(EVEN, a)))
	equal_offsetcoord("conversion_roundtrip even-q", b, QoffsetFromCube(EVEN, QoffsetToCube(EVEN, b)))
	equal_hex("conversion_roundtrip odd-q", a, QoffsetToCube(ODD, QoffsetFromCube(ODD, a)))
	equal_offsetcoord("conversion_roundtrip odd-q", b, QoffsetFromCube(ODD, QoffsetToCube(ODD, b)))
	equal_hex("conversion_roundtrip even-r", a, RoffsetToCube(EVEN, RoffsetFromCube(EVEN, a)))
	equal_offsetcoord("conversion_roundtrip even-r", b, RoffsetFromCube(EVEN, RoffsetToCube(EVEN, b)))
	equal_hex("conversion_roundtrip odd-r", a, RoffsetToCube(ODD, RoffsetFromCube(ODD, a)))
	equal_offsetcoord("conversion_roundtrip odd-r", b, RoffsetFromCube(ODD, RoffsetToCube(ODD, b)))
}

func TestOffsetFromCube(t *testing.T) {
	equal_offsetcoord("offset_from_cube even-q", NewOffsetCoord(1, 3), QoffsetFromCube(EVEN, NewHex(1, 2, -3)))
	equal_offsetcoord("offset_from_cube odd-q", NewOffsetCoord(1, 2), QoffsetFromCube(ODD, NewHex(1, 2, -3)))
}

func TestOffsetToCube(t *testing.T) {
	equal_hex("offset_to_cube even-", NewHex(1, 2, -3), QoffsetToCube(EVEN, NewOffsetCoord(1, 3)))
	equal_hex("offset_to_cube odd-q", NewHex(1, 2, -3), QoffsetToCube(ODD, NewOffsetCoord(1, 2)))
}

func TestDoubledRoundtrip(t *testing.T) {
	a := NewHex(3, 4, -7)
	b := NewDoubledCoord(1, -3)
	equal_hex("conversion_roundtrip doubled-q", a, QdoubledToCube(QdoubledFromCube(a)))
	equal_doubledcoord("conversion_roundtrip doubled-q", b, QdoubledFromCube(QdoubledToCube(b)))
	equal_hex("conversion_roundtrip doubled-r", a, RdoubledToCube(RdoubledFromCube(a)))
	equal_doubledcoord("conversion_roundtrip doubled-r", b, RdoubledFromCube(RdoubledToCube(b)))
}

func TestDoubledFromCube(t *testing.T) {
	equal_doubledcoord("doubled_from_cube doubled-q", NewDoubledCoord(1, 5), QdoubledFromCube(NewHex(1, 2, -3)))
	equal_doubledcoord("doubled_from_cube doubled-r", NewDoubledCoord(4, 2), RdoubledFromCube(NewHex(1, 2, -3)))
}

func TestDoubledToCube(t *testing.T) {
	equal_hex("doubled_to_cube doubled-q", NewHex(1, 2, -3), QdoubledToCube(NewDoubledCoord(1, 5)))
	equal_hex("doubled_to_cube doubled-r", NewHex(1, 2, -3), RdoubledToCube(NewDoubledCoord(4, 2)))
}

