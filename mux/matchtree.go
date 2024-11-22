package mux

import "net/http"

// MatchTree represents a binary tree for matching HTTP requests.
type MatchTree struct {
	matcher   func(http.Request) bool
	predicate string
	left      *MatchTree
	right     *MatchTree
}

// newMatchTree creates a MatchTree based on the given Tree structure.
func newMatchTree(t *Tree) (*MatchTree, error) {
	if t == nil {
		return nil, nil
	}

	m := &MatchTree{}

	// Set matcher function if available in funcs map.
	if f, ok := funcs[t.fn]; ok {
		if err := f(m, t.value...); err != nil {
			return nil, err
		}
	}

	// Process predicate logic and recursively create left and right subtrees.
	if t.predicate != "" {
		m.predicate = t.predicate

		// Recursively initialize the left subtree.
		leftTree, err := newMatchTree(t.left)
		if err != nil {
			return nil, err
		}
		m.left = leftTree

		// Recursively initialize the right subtree.
		rightTree, err := newMatchTree(t.right)
		if err != nil {
			return nil, err
		}
		m.right = rightTree
	}

	return m, nil
}

// Match evaluates the HTTP request against the MatchTree.
func (m *MatchTree) Match(req http.Request) bool {
	// If matcher is defined, use it.
	if m.matcher != nil {
		return m.matcher(req)
	}

	// Evaluate logical predicates if present.
	switch m.predicate {
	case "and":
		return m.left != nil && m.right != nil && m.left.Match(req) && m.right.Match(req)
	case "or":
		return (m.left != nil && m.left.Match(req)) || (m.right != nil && m.right.Match(req))
	}

	// Default to false if no matcher or predicate is defined.
	return false
}
