package main

import (
	"fmt"
)


type Tree[T constraints.Ordered] struct {
	value  T
	left   *Tree[T]
	right  *Tree[T]
	parent *Tree[T]
}


func (t *Tree[T]) Insert(v T) *Tree[T] {

	if t == nil {
		return &Tree[T]{value: v}
	}


	if v < t.value {
		child := t.left.Insert(v)
		t.left = child
		child.parent = t

	} else if v > t.value {
		child := t.right.Insert(v)
		t.right = child
		child.parent = t
	}
	// Duplicates ignored
	return t
}


func (t *Tree[T]) Search(v T) *Tree[T] {
	if t == nil {
		return nil
	}

	if v == t.value {
		return t
	} else if v < t.value {
		return t.left.Search(v)
	} else {
		return t.right.Search(v)
	}

}


func inorder[T constraints.Ordered](t *Tree[T]) {
	
	if t == nil {
		return
	}

	inorder(t.left)
	fmt.Print(t.value, " ")
	inorder(t.right)
}


func CollectInorder[T constraints.Ordered](t *Tree[T], out *[]T) {
	if t == nil {
		return
	}
	CollectInorder(t.left, out)
	*out = append(*out, t.value)
	CollectInorder(t.right, out)
}


func colCharCount (string text) []int {

	counts := make(map[rune]int)

	for _, r := range text {
		counts[r]++
	}

}


func main() {
	var root *Tree[int]

	
	for _, v := range []int{8, 3, 10, 1, 6, 14, 4, 7, 13} {
		root = root.Insert(v)
	}


	fmt.Print("In-order traversal: ")
	inorder(root)
	fmt.Println()


	if found := root.Search(6); found != nil {
		fmt.Println("Found:", found.value)
	} else {
		fmt.Println("Value 6 not found")
	}

	if found := root.Search(100); found != nil {
		fmt.Println("Found:", found.value)
	} else {
		fmt.Println("Value 100 not found")
	}

	
	var values []int
	CollectInorder(root, &values)
	fmt.Println("Collected sorted values:", values)
}