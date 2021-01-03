package itree

// Interval defines a half-open interval [Start,End).
type Interval struct {
	// Left and Right define the left and right endpoints of this interval,
	// respectively. For a non-empty interval, Left.Less(Right) holds.
	// Only non-empty intervals can be inserted into a tree.
	Left, Right Point
}

// empty reports whether this interval is empty.
func (iv Interval) empty() bool {
	return lessOrEqual(iv.Right, iv.Left)
}

// Less checks whether this interval is less than the given interval,
// by lexicographic ordering of the left and right endpoints.
func (iv Interval) Less(than Interval) bool {
	if iv.Left.Less(than.Left) {
		return true
	}
	if than.Left.Less(iv.Left) {
		return false
	}
	return iv.Right.Less(than.Right)
}

// Equal checks whether this interval is equal to the given interval.
func (iv Interval) Equal(to Interval) bool {
	return equal(iv.Left, to.Left) && equal(iv.Right, to.Right)
}

// ContainsPoint checks whether this interval contains the given point.
func (iv Interval) ContainsPoint(x Point) bool {
	return x.Less(iv.Right) && lessOrEqual(iv.Left, x)
}

// ContainsInterval checks whether this interval completely contains the
// given other interval.
func (iv Interval) ContainsInterval(other Interval) bool {
	return lessOrEqual(iv.Left, other.Left) && lessOrEqual(other.Right, iv.Right)
}

// Overlaps checks whether the intersection of this interval with the given
// interval contains at least one point.
func (iv Interval) Overlaps(with Interval) bool {
	return !(lessOrEqual(with.Right, iv.Left) || lessOrEqual(iv.Right, with.Left))
}
