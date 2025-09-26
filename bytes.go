package main

import (
	"fmt"
)

func main() {
	s := "Go"

	bytes := []byte(s)
	fmt.Println("Bytes: ", bytes)

	for _, b := range bytes {
		fmt.Printf("Byte: %08b\n", b)
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			fmt.Print(bit)
		}
		fmt.Println()
	}

}
