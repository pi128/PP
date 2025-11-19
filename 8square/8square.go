package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	list := inList()
	cube := To2D(list)

	fmt.Println("Start state (1D):", list)
	fmt.Println("Start state (2D):")
	PrintCube(cube)
	fmt.Println("Start heuristic:", heuristicSum(cube))

	goalNode := bestFirstSearch(cube)
	if goalNode == nil {
		fmt.Println("No solution found (or search limit reached).")
		return
	}

	path := []*Node{}
	for n := goalNode; n != nil; n = n.parent {
		path = append(path, n)
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	fmt.Println("\nSolution path:")
	for step, n := range path {
		fmt.Printf("Step %d, h = %d\n", step, n.h)
		PrintCube(n.cube)
		fmt.Println()
	}
	fmt.Println("Total steps:", len(path)-1)
}

type Node struct {
	cube   [][]int
	h      int
	g      int
	parent *Node
}

func bestFirstSearch(start [][]int) *Node {
	startNode := &Node{
		cube:   copyCube(start),
		h:      heuristicSum(start),
		g:      0,
		parent: nil,
	}

	frontier := []*Node{startNode}

	visited := make(map[string]bool)
	visited[cubeID(startNode.cube)] = true

	goal := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 0},
	}

	steps := 0
	maxSteps := 50000 // safety cap

	for len(frontier) > 0 && steps < maxSteps {

		bestIdx := 0
		for i := 1; i < len(frontier); i++ {
			if frontier[i].h < frontier[bestIdx].h {
				bestIdx = i
			}
		}
		current := frontier[bestIdx]

		frontier[bestIdx] = frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

		steps++

		if boardsEqual(current.cube, goal) || current.h == 0 {
			fmt.Println("Goal found in", steps, "expansions.")
			return current
		}

		neighbors := generateNeighbors(current)
		for _, nb := range neighbors {
			id := cubeID(nb.cube)
			if visited[id] {
				continue
			}
			visited[id] = true
			frontier = append(frontier, nb)
		}
	}

	return nil
}

func generateNeighbors(n *Node) []*Node {
	cube := n.cube

	var bi, bj int
	found := false
	for i := range cube {
		for j := range cube[i] {
			if cube[i][j] == 0 {
				bi, bj = i, j
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	var neighbors []*Node
	for _, d := range directions {
		ni := bi + d.di
		nj := bj + d.dj

		if ni >= 0 && ni < 3 && nj >= 0 && nj < 3 {
			newCube := copyCube(cube)

			newCube[bi][bj], newCube[ni][nj] = newCube[ni][nj], newCube[bi][bj]

			h := heuristicSum(newCube)
			nb := &Node{
				cube:   newCube,
				h:      h,
				g:      n.g + 1,
				parent: n,
			}
			neighbors = append(neighbors, nb)
		}
	}

	return neighbors
}

func copyCube(cube [][]int) [][]int {
	newCube := make([][]int, len(cube))
	for i := range cube {
		newCube[i] = make([]int, len(cube[i]))
		copy(newCube[i], cube[i])
	}
	return newCube
}

func cubeID(cube [][]int) string {
	b := make([]byte, 0, 9)
	for i := range cube {
		for j := range cube[i] {
			b = append(b, byte(cube[i][j])+'0')
		}
	}
	return string(b)
}

func boardsEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func inList() []int {
	list := []int{}

	fmt.Println("Creating a 8square, type r if want random. (Must be 0-8 and do duplicates)")
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < 9; i++ {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "r" {
			list := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
			list = randomize(list)
			return list
		}

		n, err := strconv.Atoi(input)
		if err != nil {
			panic(err)
		}

		list = append(list, n)
	}

	return list
}

func randomize(list []int) []int {
	for i := len(list) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func heuristic(cube [][]int) []int {
	goal := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 0},
	}
	mandist := []int{}

	for i := 0; i < len(cube); i++ {
		for j := 0; j < len(cube[i]); j++ {
			if cube[i][j] == goal[i][j] {
				continue
			} else {
				value := cube[i][j]
				if value == 0 {
					continue
				}

				for x := 0; x < len(cube); x++ {
					for y := 0; y < len(cube[x]); y++ {
						if goal[x][y] == value {
							mandist = append(mandist, abs(i-x)+abs(j-y))
						}
					}
				}
			}
		}
	}
	return mandist
}

func heuristicSum(cube [][]int) int {
	ds := heuristic(cube)
	total := 0
	for _, d := range ds {
		total += d
	}
	return total
}

type Delta struct {
	di int
	dj int
}

var directions = []Delta{
	{-1, 0}, // up
	{1, 0},  // down
	{0, -1}, // left
	{0, 1},  // right
}

func validmove(cube [][]int) [][]int {
	var bi, bj int
	for i := range cube {
		for j := range cube[i] {
			if cube[i][j] == 0 {
				bi, bj = i, j
			}
		}
	}

	valid := []int{}
	for _, d := range directions {
		ni := bi + d.di
		nj := bj + d.dj

		if ni >= 0 && ni < 3 && nj >= 0 && nj < 3 {
			valid = append(valid, cube[ni][nj])
		}
	}

	fmt.Println("valid moves:", valid)
	return cube
}

func To2D(list []int) [][]int {
	cube := make([][]int, 3)
	for i := range cube {
		cube[i] = make([]int, 3)
	}

	k := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cube[i][j] = list[k]
			k++
		}
	}

	return cube
}

func PrintCube(cube [][]int) {
	for i := 0; i < len(cube); i++ {
		fmt.Print("[")
		for j := 0; j < len(cube[i]); j++ {
			fmt.Print(cube[i][j])
		}
		fmt.Println("]")
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
