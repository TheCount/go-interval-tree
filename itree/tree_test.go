package itree

import (
	"math/rand"
	"testing"
)

// TestZero tests the tree zero value.
func TestZero(t *testing.T) {
	var zero T
	testInvariants(t, &zero)
	if zero.Len() != 0 {
		t.Errorf("tree zero value length %d != 0", zero.Len())
	}
	if zero.GetMin() != nil {
		t.Error("tree zero value has minimum")
	}
	if zero.GetMax() != nil {
		t.Error("tree zero value has maximum")
	}
}

// TestInsert tests inserting a value in a tree.
func TestInsert(t *testing.T) {
	var tree T
	expectPanic(t, "insert (nil, nil) interval", func() {
		tree.ReplaceOrInsert(Interval{}, 0)
	})
	expectPanic(t, "insert (x, nil) interval", func() {
		tree.ReplaceOrInsert(Interval{Int(0), nil}, 0)
	})
	expectPanic(t, "insert (nil, x) interval", func() {
		tree.ReplaceOrInsert(Interval{nil, Int(0)}, 0)
	})
	expectPanic(t, "insert empty interval", func() {
		tree.ReplaceOrInsert(Interval{Int(0), Int(0)}, 0)
	})
	interval := Interval{Int(0), Int(1)}
	previous, present := tree.ReplaceOrInsert(interval, 0)
	if present {
		t.Error("present == true after insert into empty tree")
	}
	if previous != nil {
		t.Errorf("non-nil previous %v after insert into empty tree", previous)
	}
	if tree.Len() != 1 {
		t.Errorf("expected tree len 1, got %d", tree.Len())
	}
	testInvariants(t, &tree)
	previous, present = tree.ReplaceOrInsert(interval, 1)
	if !present {
		t.Error("present == false after replace")
	}
	if previous != 0 {
		t.Errorf("previous == %v, expected 0", previous)
	}
	if tree.Len() != 1 {
		t.Errorf("expected tree len 1 after overwrite, got %d", tree.Len())
	}
	testInvariants(t, &tree)
}

// TestGet tests getting a value from a tree.
func TestGet(t *testing.T) {
	var tree T
	expectPanic(t, "get (nil, nil) interval", func() {
		tree.Get(Interval{})
	})
	expectPanic(t, "get (x, nil) interval", func() {
		tree.Get(Interval{Int(0), nil})
	})
	expectPanic(t, "get (nil, x) interval", func() {
		tree.Get(Interval{nil, Int(0)})
	})
	expectPanic(t, "get empty interval", func() {
		tree.Get(Interval{Int(0), Int(0)})
	})
	interval := Interval{Int(0), Int(1)}
	v, ok := tree.Get(interval)
	if ok {
		t.Error("ok == true after get from empty tree")
	}
	if v != nil {
		t.Errorf("got value %v from empty tree", v)
	}
	tree.ReplaceOrInsert(interval, 0)
	v, ok = tree.Get(interval)
	if !ok {
		t.Error("unable to get value from tree")
	}
	if v != 0 {
		t.Errorf("expected v == 0 after get, got %v", v)
	}
	testInvariants(t, &tree)
}

// TestGetMinMax tests getting the minimum and maximum element from a tree.
func TestGetMinMax(t *testing.T) {
	var tree T
	interval0 := Interval{Int(0), Int(1)}
	interval1 := Interval{Int(1), Int(2)}
	interval2 := Interval{Int(2), Int(3)}
	min, max := tree.GetMin(), tree.GetMax()
	if min != nil || max != nil {
		t.Error("non-nil min/max in empty tree")
	}
	tree.ReplaceOrInsert(interval1, 1)
	testInvariants(t, &tree)
	min, max = tree.GetMin(), tree.GetMax()
	if min == nil {
		t.Fatal("nil min in non-empty tree")
	}
	if min != max {
		t.Error("min != max in tree with one element")
	}
	if min.Value != 1 {
		t.Errorf("expected min/max value of 1, got %v", min.Value)
	}
	tree.ReplaceOrInsert(interval2, 2)
	testInvariants(t, &tree)
	min2, max2 := tree.GetMin(), tree.GetMax()
	if min2 != min {
		t.Error("min != old min after inserting larger interval")
	}
	if max == max2 {
		t.Error("max == old max after inserting new maximum")
	}
	if max2 == nil {
		t.Fatal("nil new max")
	}
	if max2.Value != 2 {
		t.Errorf("expected new max value of 2, got %v", max2.Value)
	}
	min, max = min2, max2
	tree.ReplaceOrInsert(interval0, 0)
	testInvariants(t, &tree)
	min2, max2 = tree.GetMin(), tree.GetMax()
	if min2 == min {
		t.Error("min == old min after inserting new minimum")
	}
	if max2 != max {
		t.Error("max != old max after inserting smaller interval")
	}
	if min2 == nil {
		t.Fatal("nil new min")
	}
	if min2.Value != 0 {
		t.Errorf("expected new min value of 0, got %v", min2.Value)
	}
}

// TestGetLessGreater tests the Get{Less,Greater}{Equal,} functions.
func TestGetLessGreater(t *testing.T) {
	var tree T
	expectPanic(t, "GetLess (nil, nil) interval", func() {
		tree.GetLess(Interval{})
	})
	expectPanic(t, "GetLess (x, nil) interval", func() {
		tree.GetLess(Interval{Int(0), nil})
	})
	expectPanic(t, "GetLess (nil, x) interval", func() {
		tree.GetLess(Interval{nil, Int(0)})
	})
	expectPanic(t, "GetLess empty interval", func() {
		tree.GetLess(Interval{Int(0), Int(0)})
	})
	expectPanic(t, "GetLessEqual (nil, nil) interval", func() {
		tree.GetLessEqual(Interval{})
	})
	expectPanic(t, "GetLessEqual (x, nil) interval", func() {
		tree.GetLessEqual(Interval{Int(0), nil})
	})
	expectPanic(t, "GetLessEqual (nil, x) interval", func() {
		tree.GetLessEqual(Interval{nil, Int(0)})
	})
	expectPanic(t, "GetLessEqual empty interval", func() {
		tree.GetLessEqual(Interval{Int(0), Int(0)})
	})
	expectPanic(t, "GetGreater (nil, nil) interval", func() {
		tree.GetGreater(Interval{})
	})
	expectPanic(t, "GetGreater (x, nil) interval", func() {
		tree.GetGreater(Interval{Int(0), nil})
	})
	expectPanic(t, "GetGreater (nil, x) interval", func() {
		tree.GetGreater(Interval{nil, Int(0)})
	})
	expectPanic(t, "GetGreater empty interval", func() {
		tree.GetGreater(Interval{Int(0), Int(0)})
	})
	expectPanic(t, "GetGreaterEqual (nil, nil) interval", func() {
		tree.GetGreaterEqual(Interval{})
	})
	expectPanic(t, "GetGreaterEqual (x, nil) interval", func() {
		tree.GetGreaterEqual(Interval{Int(0), nil})
	})
	expectPanic(t, "GetGreaterEqual (nil, x) interval", func() {
		tree.GetGreaterEqual(Interval{nil, Int(0)})
	})
	expectPanic(t, "GetGreaterEqual empty interval", func() {
		tree.GetGreaterEqual(Interval{Int(0), Int(0)})
	})
	interval0 := Interval{Int(0), Int(1)}
	interval1 := Interval{Int(1), Int(2)}
	interval2 := Interval{Int(2), Int(3)}
	l, le, g, ge := tree.GetLess(interval1), tree.GetLessEqual(interval1),
		tree.GetGreater(interval1), tree.GetGreaterEqual(interval1)
	if l != nil || le != nil || g != nil || ge != nil {
		t.Error("non-nil Get* return from empty tree")
	}
	tree.ReplaceOrInsert(interval1, 1)
	testInvariants(t, &tree)
	l, le, g, ge = tree.GetLess(interval0), tree.GetLessEqual(interval0),
		tree.GetGreater(interval0), tree.GetGreaterEqual(interval0)
	if l != nil || le != nil {
		t.Error("unexpected non-nil GetLess*(0) return from tree with one element")
	}
	if ge == nil || g == nil {
		t.Fatal("unexpected nil GetGreater*(0) from tree with one element")
	}
	if ge != g || ge.Value != 1 {
		t.Errorf("expected GetGreater* == 1, got (%v,%v)", ge.Value, g.Value)
	}
	l, le, g, ge = tree.GetLess(interval1), tree.GetLessEqual(interval1),
		tree.GetGreater(interval1), tree.GetGreaterEqual(interval1)
	if l != nil || g != nil {
		t.Error("unexpected non-nil Get*(1) return from tree with one element")
	}
	if le == nil || ge == nil {
		t.Fatal("unexpected nil Get*Equal(1) return from tree with one element")
	}
	if le != ge || le.Value != 1 {
		t.Errorf("expected Get*Equal == 1, got (%v,%v)", le.Value, ge.Value)
	}
	l, le, g, ge = tree.GetLess(interval2), tree.GetLessEqual(interval2),
		tree.GetGreater(interval2), tree.GetGreaterEqual(interval2)
	if g != nil || ge != nil {
		t.Error("unexpected non-nil GetGreater*(2) from tree with one element")
	}
	if l == nil || le == nil {
		t.Fatal("unexpected nil GetLess*(2) from tree with one element")
	}
	if l != le || le.Value != 1 {
		t.Errorf("expected GetLess* == 1, got (%v,%v)", l.Value, le.Value)
	}
	tree.ReplaceOrInsert(interval0, 0)
	tree.ReplaceOrInsert(interval2, 2)
	testInvariants(t, &tree)
	l, le, g, ge = tree.GetLess(interval0), tree.GetLessEqual(interval0),
		tree.GetGreater(interval0), tree.GetGreaterEqual(interval0)
	if l != nil {
		t.Error("unexpected non-nil GetLess(0) from tree with three elements")
	}
	if le == nil || ge == nil || g == nil {
		t.Fatal("unexpected nil GetLE/GE/G(0) from tree with three elements")
	}
	if le != ge || le.Value != 0 || g.Value != 1 {
		t.Errorf("expected GetLE/GE/G = (0,0,1), got (%v,%v,%v)",
			le.Value, ge.Value, g.Value)
	}
	l, le, g, ge = tree.GetLess(interval1), tree.GetLessEqual(interval1),
		tree.GetGreater(interval1), tree.GetGreaterEqual(interval1)
	if l == nil || le == nil || ge == nil || g == nil {
		t.Fatal("unexpected nil Get*(1) from tree with three elements")
	}
	if l.Value != 0 || le != ge || le.Value != 1 || g.Value != 2 {
		t.Errorf("expected Get* = (0,1,1,2), got (%v,%v,%v,%v)",
			l.Value, le.Value, ge.Value, g.Value)
	}
	l, le, g, ge = tree.GetLess(interval2), tree.GetLessEqual(interval2),
		tree.GetGreater(interval2), tree.GetGreaterEqual(interval2)
	if g != nil {
		t.Error("unexpected non-nil GetGreater(2) from tree with three elements")
	}
	if l == nil || le == nil || ge == nil {
		t.Fatal("unexpected nil GetL/LE/GE(2) from tree with three elements")
	}
	if l.Value != 1 || le != ge || le.Value != 2 {
		t.Errorf("expected GetL/LE/GE = (1,2,2), got (%v,%v,%v)",
			l.Value, le.Value, ge.Value)
	}
}

// TestDelete tests deleting elements from a tree.
func TestDelete(t *testing.T) {
	var tree T
	expectPanic(t, "Delete (nil, nil) interval", func() {
		tree.Delete(Interval{})
	})
	expectPanic(t, "Delete (x, nil) interval", func() {
		tree.Delete(Interval{Int(0), nil})
	})
	expectPanic(t, "Delete (nil, x) interval", func() {
		tree.Delete(Interval{nil, Int(0)})
	})
	expectPanic(t, "Delete empty interval", func() {
		tree.Delete(Interval{Int(0), Int(0)})
	})
	interval0 := Interval{Int(0), Int(1)}
	interval1 := Interval{Int(1), Int(2)}
	interval2 := Interval{Int(2), Int(3)}
	value, deleted := tree.Delete(interval0)
	if deleted {
		t.Error("got deleted == true deleting from empty tree")
	}
	if value != nil {
		t.Errorf("got previous == %v deleting from empty tree", value)
	}
	testInvariants(t, &tree)
	tree.ReplaceOrInsert(interval0, 0)
	tree.ReplaceOrInsert(interval1, 1)
	tree.ReplaceOrInsert(interval2, 2)
	testInvariants(t, &tree)
	value, deleted = tree.Delete(interval1)
	testInvariants(t, &tree)
	if !deleted {
		t.Error("failed to delete interval1")
	}
	if value != 1 {
		t.Errorf("expected deleted value to be 1, got %v", value)
	}
	value, deleted = tree.Delete(interval1)
	testInvariants(t, &tree)
	if deleted {
		t.Error("deleted == true on already deleted element")
	}
	if value != nil {
		t.Errorf("expected already deleted value == nil, got %v", value)
	}
}

// TestRandom randomly inserts and deletes elements into/from a tree.
func TestRandom(t *testing.T) {
	seedOnce.Do(seedRand)
	checkMap := make(map[Interval]interface{})
	var tree T
	for i := 0; i != 10000; i++ {
		deleteProb := float64(tree.Len()) / 1000.0
		if rand.Float64() < deleteProb {
			// delete random element
			var iv Interval
			for key := range checkMap {
				iv = key
				break
			}
			previous, deleted := tree.Delete(iv)
			testInvariants(t, &tree)
			if !deleted {
				t.Error("failed to delete random element")
			}
			if previous != checkMap[iv] {
				t.Errorf("expected deleted value = %v, got %v", checkMap[iv], previous)
			}
			delete(checkMap, iv)
		} else {
			// add random element
			iv := randomInterval()
			tree.ReplaceOrInsert(iv, i)
			testInvariants(t, &tree)
			checkMap[iv] = i
		}
	}
}
