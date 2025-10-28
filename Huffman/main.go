package main

import (
	"cmp"
	"fmt"
)

type Tree[T cmp.Ordered] struct {
	parent *Tree[T]
	left   *Tree[T]
	right  *Tree[T]
	value  T
}

// Inserts v into the BST rooted at t and returns the node that holds v.
// Works even if t is nil, as long as the caller assigns the return to root.
func (t *Tree[T]) addChild(v T) *Tree[T] {
	if t == nil {
		return &Tree[T]{value: v}
	}

	if v < t.value {
		if t.left == nil {
			t.left = &Tree[T]{value: v, parent: t}
			return t.left
		}
		return t.left.addChild(v)
	} else if v > t.value {
		if t.right == nil {
			t.right = &Tree[T]{value: v, parent: t}
			return t.right
		}
	}
	return t
}

func main() {
	fmt.Println("Hello")
}
