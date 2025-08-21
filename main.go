package main

import (
	"fmt"
)

func main() {
	fmt.Printf("hello world!\n")
	for range []int{1, 2, 3, 4} {
		fmt.Printf("hi\t")
	}

	hello := "hello"

	print(hello)

	fmt.Println("\n")
}
