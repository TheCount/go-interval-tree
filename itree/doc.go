// Package itree provides an interval tree implementation, a mutable data
// structure mapping intervals to arbitrary values.
//
// Intervals are half-open, defined by their endpoints [start,end). The tree
// can be queried for items containing a point or an interval, as well as
// various overlap conditions.
package itree
