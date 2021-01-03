package itree

import (
	"testing"
)

// blackDepth returns the black depth of the subtree rooted at n.
func blackDepth(n *Node) int {
	if n == nil {
		return 1
	}
	result := blackDepth(n.left)
	if !n.red {
		result++
	}
	return result
}

// testInvariants tests the tree invariants.
func testInvariants(t *testing.T, tree *T) {
	bd := blackDepth(tree.root)
	n := testSubtreeInvariants(t, tree.root, 0, bd, nil, nil)
	if tree.length != n {
		t.Errorf("tree expected size %d actual size %d", tree.length, n)
	}
}

// testSubtreeInvariants tests the tree invariants for the subtree rooted at n.
// currentBD and finalBD are the current black depth and the expected final
// black depth at a nil descendant, respectively; start and end define the
// acceptable range for the left interval endpoint.
// The total number of nodes in the subtree is returned.
func testSubtreeInvariants(
	t *testing.T, n *Node, currentBD, finalBD int, start, end Point,
) int {
	if n.black() {
		currentBD++
	}
	if n == nil {
		if currentBD != finalBD {
			t.Errorf("expected black depth = %d, got %d", finalBD, currentBD)
		}
		return 0
	}
	if n.red && !(n.left.black() && n.right.black()) {
		t.Errorf("red node %v with red child", n.Value)
	}
	if (start != nil && n.interval.Left.Less(start)) ||
		(end != nil && end.Less(n.interval.Left)) {
		t.Errorf("node %v -> %v out of range [%v,%v]",
			n.interval, n.Value, start, end)
	}
	expectedMaxRight := n.interval.Right
	if n.left != nil && expectedMaxRight.Less(n.left.maxRight) {
		expectedMaxRight = n.left.maxRight
	}
	if n.right != nil && expectedMaxRight.Less(n.right.maxRight) {
		expectedMaxRight = n.right.maxRight
	}
	if !equal(expectedMaxRight, n.maxRight) {
		t.Errorf("node %v expected maxRight %v, have %v",
			n.Value, expectedMaxRight, n.maxRight)
	}
	return 1 + testSubtreeInvariants(
		t, n.left, currentBD, finalBD, start, n.interval.Left,
	) + testSubtreeInvariants(
		t, n.right, currentBD, finalBD, n.interval.Left, end,
	)
}
