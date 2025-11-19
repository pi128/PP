package main

import (
	"fmt"

	"math/rand"
)

func main() {

	list := randomize()
	cube := To2D(list)

	fmt.Println(list)
	PrintCube(cube)
	move(cube)

	fmt.Println(heuristicSum(cube))
	// I can't finish this in time but we are going to learn how to do A* best first search the main idea is move the nicest looking H of the available then record the board and the next up to 4 moves again pick the best h and continue check the latest board against the previously seen one and if it has then back up and pick second best H once finished you can walk up all parent nodes to show the path

	// I also want to add in to choose the input and outputs and i bet i could vibe up a nice webpage with images and sliding

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
			//outer
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

func move(cube [][]int) [][]int {

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

func randomize() []int {

	list := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	for i := len(list) - 1; i > 0; i-- {

		j := rand.Intn(i + 1)

		list[i], list[j] = list[j], list[i]

	}

	return list
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
