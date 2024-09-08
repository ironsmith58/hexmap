package hexmap

// Adapted to golang, 2022-08-15, ironsmith58
// started with article, https://www.redblobgames.com/grids/hexagons/implementation.html
// and code from https://www.redblobgames.com/grids/hexagons/codegen/output/lib.cpp

import (
	"fmt"
	"math"
)

/*
 * Integer Absolute Value
 */
func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

/*
 * Integer Max of 2 numbers
 */
func imax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

/* PI to enough digits for our purposes */
const M_PI = 3.14159265358979323846

type Point struct {
	x float64
	y float64
}

func NewPoint(x float64, y float64) Point {
	return Point{x, y}
}

type Hex struct {
	q, r, s int
}

func NewHex(q int, r int, s int) Hex {
	if (q + r + s) != 0 {
		errstr := fmt.Sprintf("q + r + s must be 0, %d, %d, %d", q, r, s)
		panic(errstr)
	}
	return Hex{q, r, s}
}

type FractionalHex struct {
	q, r, s float64
}

func NewFractionalHex(q float64, r float64, s float64) FractionalHex {
	if math.Round(q+r+s) != 0 {
		errstr := fmt.Sprintf("q + r + s must be 0, %f, %f,%f", q, r, s)
		panic(errstr)
	}
	return FractionalHex{q, r, s}
}

type OffsetCoord struct {
	col, row int
}

func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{col, row}
}

type DoubledCoord struct {
	col, row int
}

func NewDoubledCoord(col, row int) DoubledCoord {
	return DoubledCoord{col, row}
}

type Orientation struct {
	f0          float64
	f1          float64
	f2          float64
	f3          float64
	b0          float64
	b1          float64
	b2          float64
	b3          float64
	start_angle float64
}

type Layout struct {
	orientation Orientation
	size        Point
	origin      Point
}

func NewLayout(o Orientation, s Point, c Point) Layout {
	return Layout{o, s, c}
}


func HexAdd(a Hex, b Hex) Hex {
	return NewHex(a.q+b.q, a.r+b.r, a.s+b.s)
}

func HexSubtract(a Hex, b Hex) Hex {
	return NewHex(a.q-b.q, a.r-b.r, a.s-b.s)
}

func HexScale(a Hex, k int) Hex {
	return NewHex(a.q*k, a.r*k, a.s*k)
}

func HexRotateLeft(a Hex) Hex {
	return NewHex(-a.s, -a.q, -a.r)
}

func HexRotateRight(a Hex) Hex {
	return NewHex(-a.r, -a.s, -a.q)
}

var hexDirections = []Hex{{1, 0, -1}, {1, -1, 0}, {0, -1, 1}, {-1, 0, 1}, {-1, 1, 0}, {0, 1, -1}}

func HexDirection(direction int) Hex {
	return hexDirections[direction]
}

func HexNeighbor(hex Hex, direction int) Hex {
	return HexAdd(hex, HexDirection(direction))
}

var hexDiagonals = []Hex{{2, -1, -1}, {1, -2, 1}, {-1, -1, 2}, {-2, 1, 1}, {-1, 2, -1}, {1, 1, -2}}

func HexDiagonalNeighbor(hex Hex, direction int) Hex {
	return HexAdd(hex, hexDiagonals[direction])
}

func HexLength(hex Hex) int {
	return int((iabs(hex.q) + iabs(hex.r) + iabs(hex.s)) / 2)
}

func HexDistance(a Hex, b Hex) int {
	return HexLength(HexSubtract(a, b))
}

func HexRound(h FractionalHex) Hex {
	qi := math.Round(h.q)
	ri := math.Round(h.r)
	si := math.Round(h.s)
	q_diff := math.Abs(qi - h.q)
	r_diff := math.Abs(ri - h.r)
	s_diff := math.Abs(si - h.s)
	if q_diff > r_diff && q_diff > s_diff {
		qi = -ri - si
	} else {
		if r_diff > s_diff {
			ri = -qi - si
		} else {
			si = -qi - ri
		}
	}
	return Hex{int(qi), int(ri), int(si)}
}

func HexLerp(a FractionalHex, b FractionalHex, t float64) FractionalHex {
	q := a.q*(1.0-t) + b.q*t
	r := a.r*(1.0-t) + b.r*t
	s := a.s*(1.0-t) + b.s*t
	return FractionalHex{q, r, s}
}

func HexLineDraw(a Hex, b Hex) []Hex {
	N := HexDistance(a, b)
	aq := float64(a.q)
	ar := float64(a.r)
	as := float64(a.s)
	bq := float64(b.q)
	br := float64(b.r)
	bs := float64(b.s)
	a_nudge := FractionalHex{aq + 1e-06, ar + 1e-06, as - 2e-06}
	b_nudge := FractionalHex{bq + 1e-06, br + 1e-06, bs - 2e-06}
	var results []Hex
	var step float64 = 1.0 / float64(imax(N, 1))
	for i := 0; i <= N; i++ {
		results = append(results, HexRound(HexLerp(a_nudge, b_nudge, step*float64(i))))
	}
	return results
}

const EVEN = 1
const ODD = -1

func QoffsetFromCube(offset int, h Hex) OffsetCoord {
	col := h.q
	row := h.r + int((h.q+offset*(h.q&1))/2)
	if offset != EVEN && offset != ODD {
		fmt.Println("ERROR: offset must be EVEN (+1) or ODD (-1)")
		return OffsetCoord{}
	}
	return OffsetCoord{col, row}
}

func QoffsetToCube(offset int, h OffsetCoord) Hex {
	q := h.col
	r := h.row - int((h.col+offset*(h.col&1))/2)
	s := -q - r
	if offset != EVEN && offset != ODD {
		fmt.Println("ERROR: offset must be EVEN (+1) or ODD (-1)")
		return Hex{}
	}
	return Hex{q, r, s}
}

func RoffsetFromCube(offset int, h Hex) OffsetCoord {
	col := h.q + int((h.r+offset*(h.r&1))/2)
	row := h.r
	if offset != EVEN && offset != ODD {
		fmt.Println("ERROR: offset must be EVEN (+1) or ODD (-1)")
		return OffsetCoord{}
	}
	return OffsetCoord{col, row}
}

func RoffsetToCube(offset int, h OffsetCoord) Hex {
	q := h.col - int((h.row+offset*(h.row&1))/2)
	r := h.row
	s := -q - r
	if offset != EVEN && offset != ODD {
		fmt.Println("ERROR: offset must be EVEN (+1) or ODD (-1)")
		return Hex{}
	}
	return Hex{q, r, s}
}

func QdoubledFromCube(h Hex) DoubledCoord {
	col := h.q
	row := 2*h.r + h.q
	return DoubledCoord{col, row}
}

func QdoubledToCube(h DoubledCoord) Hex {
	q := h.col
	r := int((h.row - h.col) / 2)
	s := -q - r
	return Hex{q, r, s}
}

func RdoubledFromCube(h Hex) DoubledCoord {
	col := 2*h.q + h.r
	row := h.r
	return DoubledCoord{col, row}
}

func RdoubledToCube(h DoubledCoord) Hex {
	q := int((h.col - h.row) / 2)
	r := h.row
	s := -q - r
	return Hex{q, r, s}
}

var layoutPointy = Orientation{math.Sqrt(3.0), math.Sqrt(3.0) / 2.0, 0.0, 3.0 / 2.0, math.Sqrt(3.0) / 3.0, -1.0 / 3.0, 0.0, 2.0 / 3.0, 0.5}
var layoutFlat = Orientation{3.0 / 2.0, 0.0, math.Sqrt(3.0) / 2.0, math.Sqrt(3.0), 2.0 / 3.0, 0.0, -1.0 / 3.0, math.Sqrt(3.0) / 3.0, 0.0}

func HexToPixel(layout Layout, h Hex) Point {
	M := layout.orientation
	size := layout.size
	origin := layout.origin
	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * size.x
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * size.y
	return NewPoint(x+origin.x, y+origin.y)
}

func PixelToHex(layout Layout, p Point) FractionalHex {
	M := layout.orientation
	size := layout.size
	origin := layout.origin
	pt := NewPoint((p.x-origin.x)/size.x, (p.y-origin.y)/size.y)
	q := M.b0*pt.x + M.b1*pt.y
	r := M.b2*pt.x + M.b3*pt.y
	return NewFractionalHex(q, r, -q-r)
}

func HexCornerOffset(layout Layout, corner int) Point {
	M := layout.orientation
	size := layout.size
	angle := 2.0 * M_PI * (M.start_angle - float64(corner)) / 6.0
	return NewPoint(size.x*math.Cos(angle), size.y*math.Sin(angle))
}

func PolygonCorners(layout Layout, h Hex) []Point {
	var corners []Point
	var center Point = HexToPixel(layout, h)
	for i := 0; i < 6; i++ {
		var offset Point = HexCornerOffset(layout, i)
		corners = append(corners, NewPoint(center.x+offset.x, center.y+offset.y))
	}
	return corners
}
