package main

import (
	"fmt"
)


type StringTree struct {
	Value    string
	Children []*StringTree
}


func (t *StringTree) AddChild(val string) *StringTree {
	child := &StringTree{Value: val}
	t.Children = append(t.Children, child)
	return child
}


func (t *StringTree) Print(indent string) {
	if t == nil {
		return
	}
	fmt.Println(indent + t.Value)
	for _, child := range t.Children {
		child.Print(indent + "    ") // add 4 spaces per level
	}
}


func (t *StringTree) Find(name string) *StringTree {
	if t == nil {
		return nil
	}
	queue := []*StringTree{t}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.Value == name {
			return cur
		}

		queue = append(queue, cur.Children...)
	}
	return nil
}


func (t *StringTree) DFS(visit func(string)) {
	if t == nil {
		return
	}
	visit(t.Value)
	for _, child := range t.Children {
		child.DFS(visit)
	}
}


func (t *StringTree) BFS(visit func(string)) {
	if t == nil {
		return
	}
	queue := []*StringTree{t}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visit(cur.Value)
		queue = append(queue, cur.Children...)
	}
}


func DemoFamilyTree() {
	root := &StringTree{Value: "Family"}
	mary := root.AddChild("Mary")
	barry := root.AddChild("Barry")
	mary.AddChild("Jock")
	mary.AddChild("Matt")
	barry.AddChild("Squid")

	fmt.Println("Family Tree:")
	root.Print("")

	fmt.Println("\nDFS:")
	root.DFS(func(v string) { fmt.Println(v) })

	fmt.Println("\nBFS:")
	root.BFS(func(v string) { fmt.Println(v) })

	fmt.Println("\nSearching for 'Squid'...")
	found := root.Find("Squid")
	if found != nil {
		fmt.Println("Found node:", found.Value)
	} else {
		fmt.Println("Not found.")
	}
}


func main() {
	fmt.Println("Hey")
	DemoFamilyTree()
}