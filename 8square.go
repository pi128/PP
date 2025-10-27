package main

import (
	"fmt"

	"math/rand"
)

func main() {
	cube := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	for i := len(cube) - 1; i > 0; i-- {

		j := rand.Intn(i + 1)

		cube[i], cube[j] = cube[j], cube[i]

	}
	fmt.Println(cube)

}
