package itree

// Node represents an element of an interval tree.
type Node struct {
	// Value is the node value. The tree structure is independent of the value.
	Value interface{}

	// interval is the interval which is mapped to Value by this node.
	interval Interval

	// maxRight is the maximal right endpoint in the subtree defined by this node,
	// i. e., maxRight = max(interval.Right, left.maxRight, right.maxRight).
	maxRight Point

	// parent, left, and right are the parent node and the left and right child,
	// respectively, of this node. May be nil. The parent is nil only if this is
	// the root node of the tree. The subtree defined by left only contains
	// nodes with intervals less than interval. The subtree defined by right
	// only contains nodes for which intervall is less than the respective node
	// interval.
	parent, left, right *Node

	// red indicates whether this node is red in the underlying red-black tree
	// structure. A black node is a nil node or a node for which red == false.
	// The following invariants hold:
	// a red node has no red children; a root node (parent == nil) is black;
	// the path from a root node to any nil descendant contains a constant number
	// of black nodes.
	red bool
}

// black reports whether the given node is black. Nil nodes are considered
// to be black.
func (n *Node) black() bool {
	return n == nil || !n.red
}

// sibling returns the sibling of this node. The given parent must be the parent
// of n. This extra parameter allows the sibling method to be called on a
// nil child. If this node does not have a sibling, nil is returned.
func (n *Node) sibling(parent *Node) *Node {
	switch {
	case parent == nil:
		return nil
	case parent.left == n:
		return parent.right
	default:
		return parent.left
	}
}

// replaceOrInsert adds the given interval → value mapping to the subtree
// defined by this node. If the interval already exists in the subtree, the
// previous value is returned with present == true. Otherwise, (nil, false) is
// returned.
func (n *Node) replaceOrInsert(t *T, interval Interval, value interface{}) (
	previous interface{}, present bool,
) {
	if n.maxRight.Less(interval.Right) {
		n.maxRight = interval.Right
	}
	switch {
	case n.interval.Less(interval):
		if n.right != nil {
			return n.right.replaceOrInsert(t, interval, value)
		}
		n.right = &Node{
			Value:    value,
			interval: interval,
			maxRight: interval.Right,
			parent:   n,
			red:      true,
		}
		t.rebalanceRed(n.right)
		t.length++
		return nil, false
	case interval.Less(n.interval):
		if n.left != nil {
			return n.left.replaceOrInsert(t, interval, value)
		}
		n.left = &Node{
			Value:    value,
			interval: interval,
			maxRight: interval.Right,
			parent:   n,
			red:      true,
		}
		t.rebalanceRed(n.left)
		t.length++
		return nil, false
	default:
		previous = n.Value
		n.Value = value
		return previous, true
	}
}

// nodesContainingPoint appends all nodes in the subtree defined by this node
// whose interval contains the given point to the given list and returns it.
func (n *Node) nodesContainingPoint(list []*Node, p Point) []*Node {
	if lessOrEqual(n.maxRight, p) {
		return list
	}
	if n.left != nil {
		list = n.left.nodesContainingPoint(list, p)
	}
	if lessOrEqual(n.interval.Left, p) {
		if p.Less(n.interval.Right) {
			list = append(list, n)
		}
		if n.right != nil {
			list = n.right.nodesContainingPoint(list, p)
		}
	}
	return list
}

// nodesContainingInterval appends all nodes in the subtree defined by this node
// whose interval contains the given interval to the given list and returns it.
func (n *Node) nodesContainingInterval(list []*Node, iv Interval) []*Node {
	if n.maxRight.Less(iv.Right) {
		return list
	}
	if n.left != nil {
		list = n.left.nodesContainingInterval(list, iv)
	}
	if lessOrEqual(n.interval.Left, iv.Left) {
		if lessOrEqual(iv.Right, n.interval.Right) {
			list = append(list, n)
		}
		if n.right != nil {
			list = n.right.nodesContainingInterval(list, iv)
		}
	}
	return list
}

// nodesContainedInInterval appends all nodes in the subtree defined by this
// node whose interval is contained in the given interval to the given list and
// returns it.
func (n *Node) nodesContainedInInterval(list []*Node, iv Interval) []*Node {
	if lessOrEqual(n.maxRight, iv.Left) {
		return list
	}
	if lessOrEqual(iv.Left, n.interval.Left) {
		if n.left != nil {
			list = n.left.nodesContainedInInterval(list, iv)
		}
		if lessOrEqual(n.interval.Right, iv.Right) {
			list = append(list, n)
		}
	}
	if n.right != nil {
		list = n.right.nodesContainedInInterval(list, iv)
	}
	return list
}

// nodesOverlappingInterval appends all nodes in the subtree defined by this
// node whose interval has a non-empty intersection with the given interval
// to the given list and returns it.
func (n *Node) nodesOverlappingInterval(list []*Node, iv Interval) []*Node {
	if lessOrEqual(n.maxRight, iv.Left) {
		return list
	}
	if iv.Left.Less(n.interval.Right) {
		if n.left != nil {
			list = n.left.nodesOverlappingInterval(list, iv)
		}
		if n.interval.Left.Less(iv.Right) {
			list = append(list, n)
		}
	}
	if n.right != nil {
		list = n.right.nodesOverlappingInterval(list, iv)
	}
	return list
}

// Interval returns a shallow copy of the interval of this node. The caller must
// not make deep changes to the returned interval which affect the result of
// Less (e. g., if the dynamic type of the interval points is a pointer type).
func (n *Node) Interval() Interval {
	return n.interval
}

// Next returns the next node in the tree sort-order (see Interval.Less).
// If no such node exists, nil is returned.
func (n *Node) Next() *Node {
	candidate := n
	if candidate.right == nil {
		for {
			if candidate.parent == nil {
				return nil
			}
			if candidate.parent.left == candidate {
				return candidate.parent
			}
			candidate = candidate.parent
		}
	}
	candidate = candidate.right
	for candidate.left != nil {
		candidate = candidate.left
	}
	return candidate
}

// Previous returns the previous node in the tree sort-order (see
// Interval.Less).
// If no such node exists, nil is returned.
func (n *Node) Previous() *Node {
	candidate := n
	if candidate.left == nil {
		for {
			if candidate.parent == nil {
				return nil
			}
			if candidate.parent.right == candidate {
				return candidate.parent
			}
			candidate = candidate.parent
		}
	}
	candidate = candidate.left
	for candidate.right != nil {
		candidate = candidate.right
	}
	return candidate
}
