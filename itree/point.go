package itree

import (
	"bytes"
	"math"
)

// Point represents a point in an interval.
type Point interface {
	// Less tests whether this point is less than the given argument.
	//
	// Less must provide a total or a strict weak ordering.
	// If the ordering is not total, for each class of incomparable points
	// (x, y are incomparable if !x.Less(y) and !y.Less(x)),
	// at most one unspecified representative will be in a given tree.
	Less(than Point) bool
}

// equal tests whether the given points are equal.
func equal(x, y Point) bool {
	return !(x.Less(y) || y.Less(x))
}

// lessOrEqual tests whether the given point x is less than or equal to the
// given point y.
func lessOrEqual(x, y Point) bool {
	return !y.Less(x)
}

// Int implements Point for the built-in int type.
type Int int

// Less checks whether x < y. It panics if y is not Int.
func (x Int) Less(y Point) bool {
	return x < y.(Int)
}

// Int8 implements Point for the built-in int8 type.
type Int8 int8

// Less checks whether x < y. It panics if y is not Int8.
func (x Int8) Less(y Point) bool {
	return x < y.(Int8)
}

// Int16 implements Point for the built-in int16 type.
type Int16 int16

// Less checks whether x < y. It panics if y is not Int16.
func (x Int16) Less(y Point) bool {
	return x < y.(Int16)
}

// Int32 implements Point for the built-in int32 type.
type Int32 int32

// Less checks whether x < y. It panics if y is not Int32.
func (x Int32) Less(y Point) bool {
	return x < y.(Int32)
}

// Int64 implements Point for the built-in int64 type.
type Int64 int64

// Less checks whether x < y. It panics if y is not Int64.
func (x Int64) Less(y Point) bool {
	return x < y.(Int64)
}

// Uint implements Point for the built-in uint type.
type Uint uint

// Less checks whether x < y. It panics if y is not Uint.
func (x Uint) Less(y Point) bool {
	return x < y.(Uint)
}

// Uint8 implements Point for the built-in uint8 type.
type Uint8 uint8

// Less checks whether x < y. It panics if y is not Uint8.
func (x Uint8) Less(y Point) bool {
	return x < y.(Uint8)
}

// Uint16 implements Point for the built-in uint16 type.
type Uint16 uint16

// Less checks whether x < y. It panics if y is not Uint16.
func (x Uint16) Less(y Point) bool {
	return x < y.(Uint16)
}

// Uint32 implements Point for the built-in uint32 type.
type Uint32 uint32

// Less checks whether x < y. It panics if y is not Uint32.
func (x Uint32) Less(y Point) bool {
	return x < y.(Uint32)
}

// Uint64 implements Point for the built-in uint64 type.
type Uint64 uint64

// Less checks whether x < y. It panics if y is not Uint64.
func (x Uint64) Less(y Point) bool {
	return x < y.(Uint64)
}

// Uintptr implements Point for the built-in uintptr type.
type Uintptr uintptr

// Less checks whether x < y. It panics if y is not Uintptr.
func (x Uintptr) Less(y Point) bool {
	return x < y.(Uintptr)
}

// Float32 implements Point for the built-in float32 type.
// NaN-values are not permitted. +0 and -0 are considered to be equal.
type Float32 float32

// Less checks whether x < y. It panics if y is not Float32 or either of x and
// y is NaN.
func (x Float32) Less(y Point) bool {
	yf := y.(Float32)
	if math.IsNaN(float64(x)) || math.IsNaN(float64(yf)) {
		panic("comparing NaN")
	}
	return x < yf
}

// Float64 implements Point for the built-in float64 type.
// NaN-values are not permitted. +0 and -0 are considered to be equal.
type Float64 float64

// Less checks whether x < y. It panics if y is not Float32 or either of x and
// y is NaN.
func (x Float64) Less(y Point) bool {
	yf := y.(Float64)
	if math.IsNaN(float64(x)) || math.IsNaN(float64(yf)) {
		panic("comparing NaN")
	}
	return x < yf
}

// String implements Point for the built-in string type.
type String string

// Less checks whether x < y. It panics if y is not String.
func (x String) Less(y Point) bool {
	return x < y.(String)
}

// Bytes implements Point for the common []byte type.
// A nil value and an empty byte slice are considered to be equal.
type Bytes []byte

// Less checks whether x < y using the bytes.Compare function.
// It panics if y is not Bytes.
func (x Bytes) Less(y Point) bool {
	return bytes.Compare(x, y.(Bytes)) < 0
}

// Point interface checks.
var (
	_ = []Point{
		Int(0), Int8(0), Int16(0), Int32(0), Int64(0),
		Uint(0), Uint8(0), Uint16(0), Uint32(0), Uint64(0), Uintptr(0),
		Float32(0), Float64(0),
		String(""), Bytes(nil),
	}
)
