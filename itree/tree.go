package itree

// T represents an interval tree.
// The zero value represents an empty tree.
// This data structure is not safe for concurrent modification.
type T struct {
	// root is the root node of this interval tree.
	root *Node

	// length is the total number of nodes in this tree.
	length int
}

// Len returns the number of elements in this interval tree.
func (t *T) Len() int {
	return t.length
}

// rotateLeft performs a left rotation of the given node m, which must have a
// right child n, and fixes the maxRight values of the affected nodes.
// Specifically, the following operation is performed (y may be nil):
//
//     m                n
//      \              /
//       n    --->    m
//      /              \
//     y                y
//
// The order defined by Interval.Less is left invariant by this operation.
func (t *T) rotateLeft(m *Node) {
	parent := m.parent
	n := m.right
	y := n.left

	// Update parents
	n.parent = parent
	m.parent = n
	if y != nil {
		y.parent = m
	}

	// Update children
	if parent != nil {
		if parent.left == m {
			parent.left = n
		} else {
			parent.right = n
		}
	} else {
		t.root = n
	}
	n.left = m
	m.right = y

	// Update maxRight
	m.maxRight = m.interval.Right
	if m.left != nil && m.maxRight.Less(m.left.maxRight) {
		m.maxRight = m.left.maxRight
	}
	if y != nil && m.maxRight.Less(y.maxRight) {
		m.maxRight = y.maxRight
	}
	if n.maxRight.Less(m.maxRight) {
		n.maxRight = m.maxRight
	}
}

// rotateRight performs a right rotation of the given node n, which must have a
// left child m, and fixes the maxRight values of the affected nodes.
// Specifically, the following operation is performed (y may be nil):
//
//     m                n
//      \              /
//       n    <---    m
//      /              \
//     y                y
//
// The order defined by Interval.Less is left invariant by this operation.
func (t *T) rotateRight(n *Node) {
	parent := n.parent
	m := n.left
	y := m.right

	// Update parents
	m.parent = parent
	n.parent = m
	if y != nil {
		y.parent = n
	}

	// Update children
	if parent != nil {
		if parent.left == n {
			parent.left = m
		} else {
			parent.right = m
		}
	} else {
		t.root = m
	}
	n.left = y
	m.right = n

	// Update maxRight
	n.maxRight = n.interval.Right
	if y != nil && n.maxRight.Less(y.maxRight) {
		n.maxRight = y.maxRight
	}
	if n.right != nil && n.maxRight.Less(n.right.maxRight) {
		n.maxRight = n.right.maxRight
	}
	if m.maxRight.Less(n.maxRight) {
		m.maxRight = n.maxRight
	}
}

// GetNode retrieves the node for the given interval from this tree.
// If no such node exists, nil is returned.
func (t *T) GetNode(iv Interval) *Node {
	if iv.empty() {
		panic("empty interval")
	}
	for current := t.root; current != nil; {
		switch {
		case iv.Less(current.interval):
			current = current.left
		case current.interval.Less(iv):
			current = current.right
		default:
			return current
		}
	}
	return nil
}

// GetMin returns the node with the lowest-sorting interval (see Interval.Less)
// in this tree. If the tree is empty, nil is returned.
func (t *T) GetMin() *Node {
	var candidate *Node
	for current := t.root; current != nil; current = current.left {
		candidate = current
	}
	return candidate
}

// GetMax returns the node with the highest-sorting interval (see Interval.Less)
// in this tree. If the tree is empty, nil is returned.
func (t *T) GetMax() *Node {
	var candidate *Node
	for current := t.root; current != nil; current = current.right {
		candidate = current
	}
	return candidate
}

// GetLess returns the node with the highest-sorting interval less than the
// given interval, or nil if no such node exists in this tree.
func (t *T) GetLess(iv Interval) *Node {
	if iv.empty() {
		panic("empty interval")
	}
	var candidate *Node
	for current := t.root; current != nil; {
		switch {
		default:
			current = current.left
		case current.interval.Less(iv):
			candidate = current
			current = current.right
		}
	}
	return candidate
}

// GetLessEqual returns the node with the highest-sorting interval less than
// or equal to the given interval, or nil if no such node exists in this tree.
func (t *T) GetLessEqual(iv Interval) *Node {
	if iv.empty() {
		panic("empty interval")
	}
	var candidate *Node
	for current := t.root; current != nil; {
		switch {
		case iv.Less(current.interval):
			current = current.left
		default:
			candidate = current
			current = current.right
		}
	}
	return candidate
}

// GetGreater returns the node with the lowest-sorting interval greater than
// the given interval, or nil if no such node exists in this tree.
func (t *T) GetGreater(iv Interval) *Node {
	if iv.empty() {
		panic("empty interval")
	}
	var candidate *Node
	for current := t.root; current != nil; {
		switch {
		case iv.Less(current.interval):
			candidate = current
			current = current.left
		default:
			current = current.right
		}
	}
	return candidate
}

// GetGreaterEqual returns the node with the lowest-sorting interval greater
// than or equal to the given interval, or nil if no such node exists in this
// tree.
func (t *T) GetGreaterEqual(iv Interval) *Node {
	if iv.empty() {
		panic("empty interval")
	}
	var candidate *Node
	for current := t.root; current != nil; {
		switch {
		default:
			candidate = current
			current = current.left
		case current.interval.Less(iv):
			current = current.right
		}
	}
	return candidate
}

// Get retrieves the value for the specified interval. If the given interval
// is not part of this tree, (nil, false) is returned. Otherwise, the value and
// present == true is returned.
func (t *T) Get(iv Interval) (value interface{}, present bool) {
	node := t.GetNode(iv)
	if node == nil {
		return nil, false
	}
	return node.Value, true
}

// ReplaceOrInsert adds the given interval → value mapping to this tree. If the
// interval already exists in the tree, the previous value is returned with
// present == true. Otherwise, (nil, false) is returned.
func (t *T) ReplaceOrInsert(interval Interval, value interface{}) (
	previous interface{}, present bool,
) {
	if interval.empty() {
		panic("empty interval")
	}
	if t.length == 0 { // empty tree
		t.root = &Node{
			Value:    value,
			interval: interval,
			maxRight: interval.Right,
		}
		t.length = 1
		return nil, false
	}

	return t.root.replaceOrInsert(t, interval, value)
}

// NodesContainingPoint returns all nodes containing the given point, ordered
// by their intervals.
// This operation runs in O(s+log(n)) time, where s is the number of returned
// nodes and n is the size of this tree.
func (t *T) NodesContainingPoint(p Point) []*Node {
	if t.root == nil {
		return nil
	}
	result := make([]*Node, 0)
	return t.root.nodesContainingPoint(result, p)
}

// NodesContainingInterval returns all nodes containing the given interval,
// ordered by their intervals.
// This operation runs in O(s+log(n)) time, where s is the number of returned
// nodes and n is the size of this tree.
func (t *T) NodesContainingInterval(iv Interval) []*Node {
	if iv.empty() {
		panic("empty interval")
	}
	if t.root == nil {
		return nil
	}
	result := make([]*Node, 0)
	return t.root.nodesContainingInterval(result, iv)
}

// NodesContainedInInterval returns all nodes contained in the given interval,
// ordered by their intervals.
// This operation runs in O(s+log(n)) time, where s is the number of returned
// nodes and n is the size of this tree.
func (t *T) NodesContainedInInterval(iv Interval) []*Node {
	if iv.empty() {
		panic("empty interval")
	}
	if t.root == nil {
		return nil
	}
	result := make([]*Node, 0)
	return t.root.nodesContainedInInterval(result, iv)
}

// NodesOverlappingInterval returns all nodes overlapping with the given
// interval, i. e., where the intersection of the node interval with iv is not
// empty.
// This operation runs in O(s+log(n)) time, where s is the number of returned
// nodes and n is the size of this tree.
func (t *T) NodesOverlappingInterval(iv Interval) []*Node {
	if iv.empty() {
		panic("empty interval")
	}
	if t.root == nil {
		return nil
	}
	result := make([]*Node, 0)
	return t.root.nodesOverlappingInterval(result, iv)
}

// DeleteNode deletes the given node from this tree. The given node must be
// part of this tree. After deletion the given node should no longer be used.
func (t *T) DeleteNode(n *Node) {
	if n.left != nil && n.right != nil {
		// n has two children, so we reduce to the one child case first by swapping
		// n with the maximum lower node in its subtree.
		// There is no need to update maxRight at this step, since maxRight will
		// have to be fixed in all ancestors of n later on anyway.
		parent := n.parent
		candidate := n.left
		if candidate.right == nil {
			// OK, we can swap n with candidate
			if parent != nil {
				if parent.left == n {
					parent.left = candidate
				} else {
					parent.right = candidate
				}
			} else {
				t.root = candidate
			}
			n.right.parent = candidate
			if candidate.left != nil {
				candidate.left.parent = n
			}
			candidate.left, n.left = n, candidate.left
			candidate.right, n.right = n.right, nil
			candidate.parent, n.parent = n.parent, candidate
		} else {
			// Search for rightmost descendant of candidate and swap with n
			for dowhile := true; dowhile; dowhile = candidate.right != nil {
				candidate = candidate.right
			}
			if parent != nil {
				if parent.left == n {
					parent.left = candidate
				} else {
					parent.right = candidate
				}
			} else {
				t.root = candidate
			}
			n.left.parent = candidate
			n.right.parent = candidate
			candidate.parent.right = n
			if candidate.left != nil {
				candidate.left.parent = n
			}
			candidate.left, n.left = n.left, candidate.left
			candidate.right, n.right = n.right, candidate.right
			candidate.parent, n.parent = n.parent, candidate.parent
		}
		n.red, candidate.red = candidate.red, n.red
	}

	// At this point, n has at most one child, so we just have to unlink n from
	// the tree and link the parent of n with the child.
	parent := n.parent
	child := n.left
	if child == nil {
		child = n.right
	}
	if child != nil {
		child.parent = parent
	}
	if parent != nil {
		if parent.left == n {
			parent.left = child
		} else {
			parent.right = child
		}
	} else {
		t.root = child
	}

	// Fix maxRight in all ancestors of n
	for current := parent; current != nil; current = current.parent {
		current.maxRight = current.interval.Right
		if current.left != nil && current.maxRight.Less(current.left.maxRight) {
			current.maxRight = current.left.maxRight
		}
		if current.right != nil && current.maxRight.Less(current.right.maxRight) {
			current.maxRight = current.right.maxRight
		}
	}

	// Finally, rebalance the tree if necessary
	if !n.red {
		if child != nil && child.red {
			child.red = false
		} else {
			t.rebalanceBlack(parent, child)
		}
	}

	t.length--
}

// DeleteAndAscend deletes the given node and returns the next node in the
// sort-order (see Interval.Less). If there is no next node, nil is returned.
// The given node should not be used afterwards.
func (t *T) DeleteAndAscend(n *Node) *Node {
	result := n.Next()
	t.DeleteNode(n)
	return result
}

// DeleteAndDescend deletes the given node and returns the previous node in the
// sort-order (see Interval.Less). If there is no previous node, nil is
// returned.
// The given node should not be used afterwards.
func (t *T) DeleteAndDescend(n *Node) *Node {
	result := n.Previous()
	t.DeleteNode(n)
	return result
}

// Delete deletes the node for the specified interval from the tree.
// If no such node exists, (nil, false) is returned. Otherwise, the value of
// the deleted node is returned with deleted == true.
func (t *T) Delete(interval Interval) (value interface{}, deleted bool) {
	n := t.GetNode(interval)
	if n == nil {
		return nil, false
	}
	value = n.Value
	t.DeleteNode(n)
	return value, true
}

// rebalanceRed rebalances the tree for the case that n is a red node whose
// parent is also red. Otherwise, rebalanceRed does nothing.
func (t *T) rebalanceRed(n *Node) {
	if n.black() { // nothing to do
		return
	}

	// If n is the root node, we can safely paint it black.
	parent := n.parent
	if parent == nil {
		n.red = false
		return
	}

	// Nothing to do if parent is not red.
	if !parent.red {
		return
	}

	// Actually rebalance
	grandparent := parent.parent // non-nil since parent is red
	auncle := parent.sibling(grandparent)
	if auncle != nil && auncle.red {
		// We have a red auncle, recursively repaint the tree up
		auncle.red = false
		parent.red = false
		grandparent.red = true
		t.rebalanceRed(grandparent)
		return
	}
	// We have a black auncle, so we rotate such that the grandparent becomes
	// n's red sibling, and n's parent becomes black.
	switch {
	case n == parent.right && parent == grandparent.left:
		t.rotateLeft(parent)
		n = n.left
		parent = grandparent.left
	case n == parent.left && parent == grandparent.right:
		t.rotateRight(parent)
		n = n.right
		parent = grandparent.right
	}
	parent.red = false
	grandparent.red = true
	if n == parent.left {
		t.rotateRight(grandparent)
	} else {
		t.rotateLeft(grandparent)
	}
}

// rebalanceBlack rebalances the tree after a black node with a black parent
// was removed. The given node n is the new child of parent (n may be nil,
// so parent must also be given). Note that in recursive calls, parent itself
// may be red.
func (t *T) rebalanceBlack(parent, n *Node) {
	if parent == nil {
		return
	}
	sibling := n.sibling(parent) // non-nil because removed node was black

	// If sibling is red, paint parent red and make sibling our black
	// grandparent through rotation.
	if sibling.red {
		parent.red = true
		sibling.red = false
		if n == parent.left {
			t.rotateLeft(parent)
		} else {
			t.rotateRight(parent)
		}
	}

	// Case 1: (new) sibling and niecews are black: fix parent recursively.
	sibling = n.sibling(parent) // may have changed due to rotation
	lniecew, rniecew := sibling.left, sibling.right
	if sibling.black() && lniecew.black() && rniecew.black() {
		sibling.red = true
		if parent.red {
			parent.red = false
		} else {
			t.rebalanceBlack(parent.parent, parent)
		}
		return
	}

	// Case 2: one of the niecews is red, rotate it up and move the paint to the
	// sibling, reducing to case 3.
	switch {
	case n == parent.left && rniecew.black(): // lniecew is red
		sibling.red = true
		lniecew.red = false
		t.rotateRight(sibling)
	case n == parent.right && lniecew.black(): // rniecew is red
		sibling.red = true
		rniecew.red = false
		t.rotateLeft(sibling)
	}

	// Case 3
	sibling = n.sibling(parent) // may have changed due to rotation
	lniecew, rniecew = sibling.left, sibling.right
	sibling.red = parent.red
	parent.red = false
	if n == parent.left {
		if rniecew != nil {
			rniecew.red = false
		}
		t.rotateLeft(parent)
	} else {
		if lniecew != nil {
			lniecew.red = false
		}
		t.rotateRight(parent)
	}
}
